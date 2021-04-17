package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	DB = "Mongo"
)

var globalRouter *Router
var globalTaskProcessor *TaskProcessor
var testEnv = flag.String("env", "local", "testEnv")
var hostUrl = flag.String("host", "127.0.0.1:8083", "host url")
var scheme = "http"

func TestMain(m *testing.M) {
	if !flag.Parsed() {
		flag.Parse()
	}
	_ = flag.Args() // flag.Args() 返回 -args 后面的所有参数，以切片表示，每个元素代表一个参数
	genTestFile()
	if *testEnv == "local" {
		logrus.Warning("test env: Local")
		initRouterAndProcessor()
	} else if *testEnv == "cloud" {
		logrus.Warning("test env: cloud")
		sendDataAndRecord("GET", "/debug/unlock_test_user", nil)
		sendDataAndRecord("GET", "/debug/drop_task_table", nil)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func genTestFile() {
	os.MkdirAll("../test/tmp/", os.ModePerm)
	buf := new(bytes.Buffer)
	alpha := image.NewAlpha(image.Rect(0, 0, 1000, 1000))
	for x := 0; x < 1000; x++ {
		for y := 0; y < 1000; y++ {
			alpha.Set(x, y, color.Alpha{uint8(x % 256)}) //设定alpha图片的透明度
		}
	}
	jpeg.Encode(buf, alpha, nil)
	err := ioutil.WriteFile("../test/tmp/test.jpeg", buf.Bytes(), 0644)
	content := []byte("jcsPan transporter Test SincereXIA @ " + time.Now().String())
	err = ioutil.WriteFile("../test/tmp/test.txt", content, 0644)
	if err != nil {
		panic(err)
	}
}

func initRouterAndProcessor() (*Router, *TaskProcessor) {
	if globalRouter != nil && globalTaskProcessor != nil {
		return globalRouter, globalTaskProcessor
	}
	var storage model.TaskStorage
	var clientDatabase model.CloudDatabase
	var fileDatabase model.FileDatabase
	//util.ReadConfigFromFile("../transporter_config.json")
	err := util.CheckConfig()
	if err != nil {
		return nil, nil
	}
	if util.Config.Database.Driver == util.MongoDB {
		util.ClearAll()
		storage, _ = model.NewMongoTaskStorage()
		clientDatabase, _ = model.NewMongoCloudDatabase()
		fileDatabase, _ = model.NewMongoFileDatabase()
	} else {
		storage = model.NewInMemoryTaskStorage()
		clientDatabase = model.NewSimpleInMemoryStorageDatabase()
		fileDatabase = model.NewInMemoryFileDatabase()
	}
	processor := TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(clientDatabase)
	// 初始化 FileInfo 数据库
	processor.FileDatabase = fileDatabase
	// 初始化 Lock
	lock, _ := NewLock(util.Config.ZookeeperHost)
	processor.Lock = lock
	processor.Lock.UnLockAll("/tester")
	// 初始化 Scheduler
	scheduler := JcsPanScheduler{
		LocalCloudID:     "aliyun-hangzhou",
		SchedulerHostUrl: "http://192.168.105.13:8082",
		ReloadCloudInfo:  true,
		CloudDatabase:    clientDatabase,
	}
	processor.Scheduler = &scheduler
	// 初始化 Monitor
	userDB, err := model.NewMongoUserDatabase()
	processor.Monitor = NewTrafficMonitor(userDB)
	// 初始化路由
	router := NewTestRouter(processor)
	// 启动 processor
	processor.StartProcessTasks(context.Background())
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	globalRouter = router
	globalTaskProcessor = &processor
	return router, &processor
}

func sendDataAndRecord(method string, url string, data io.Reader) *http.Response {
	if *testEnv == "local" {
		req, _ := http.NewRequest(method, url, data)
		recorder := httptest.NewRecorder()
		globalRouter.ServeHTTP(recorder, req)
		return recorder.Result()
	} else if *testEnv == "cloud" {
		req, _ := http.NewRequest(method, scheme+"://"+*hostUrl+url, data)
		resp, _ := http.DefaultClient.Do(req)
		return resp
	}
	return nil
}

func sendRequestAndRecord(req *http.Request) (resp *http.Response) {
	if *testEnv == "local" {
		recorder := httptest.NewRecorder()
		globalRouter.ServeHTTP(recorder, req)
		return recorder.Result()
	} else if *testEnv == "cloud" {
		req.Host = *hostUrl
		req.URL.Scheme = scheme
		req.URL.Host = *hostUrl
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Errorf("send request err: %v", err)
		}
		return resp
	}
	return nil
}

func setCookie(req *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "sid",
		Value:   "tttteeeesssstttt",
		Expires: expire,
	}
	req.AddCookie(&cookie)
}

func waitProcessorAllDone() {
	for true {
		time.Sleep(time.Millisecond * 500)
		if *testEnv == "local" {
			if globalTaskProcessor.taskStorage.IsAllDone() {
				return
			}
		} else if *testEnv == "cloud" {
			resp := sendDataAndRecord("GET", "/state/process_state", nil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			result := buf.String()
			if result == "done" {
				return
			}
		}
	}
}

func getDownloadUrl(fileID string) string {
	if *testEnv == "local" {
		fileInfo, err := globalTaskProcessor.FileDatabase.GetFileInfo(fileID)
		if err != nil {
			logrus.Fatalf("get file info err:%v", err)
		}
		url := fileInfo.DownloadUrl
		if url == "" {
			logrus.Fatalf("get download url err")
		}
		return url
	} else if *testEnv == "cloud" {
		resp := sendDataAndRecord("GET", "/debug/get_file_download_url?id="+fileID, nil)
		if resp.StatusCode != http.StatusOK {
			return ""
		}
		bodyRes, _ := ioutil.ReadAll(resp.Body)
		result := string(bodyRes)
		return result
	}
	return ""
}

func TestECUploadAndDownload(t *testing.T) {
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`
{
  "TaskType": "Upload",
   "UserID": "tester",
   "DestinationPath":"path/to/jcspantest.txt",
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ],
      "N": 3,
      "K": 2
   }
}`)
		resp := sendDataAndRecord("POST", "/task", bytes.NewBuffer(jsonStr))
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		if reply.Code != http.StatusOK {
			t.Fatalf("create upload task fail: %v", reply.Msg)
		}

		token := reply.Data.Result

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()
		req, _ := postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", token)
		setCookie(req)
		sendRequestAndRecord(req)
		waitProcessorAllDone()
	})
	t.Run("Create EC Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "UserID": "tester",
     "Sid": "tttteeeesssstttt",
     "SourcePath":"path/to/jcspantest.txt",
     "SourceStoragePlan":{
        "StorageMode": "EC",
        "Clouds": [
           {
              "CloudID": "aliyun-beijing"
           },
           {
              "CloudID": "aliyun-beijing"
           },
           {
              "CloudID": "aliyun-beijing"
           }
        ],
        "N": 3,
        "K": 2
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		tid := reply.Data.Result
		logrus.Debugf("tid: %v", tid)
		waitProcessorAllDone()
	})
	var url string
	t.Run("Check File DB and get download url", func(t *testing.T) {
		url = getDownloadUrl("tester/path/to/jcspantest.txt")
		t.Logf("download url: %v", url)
		waitProcessorAllDone()
	})
	t.Run("Get file", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		resp := sendRequestAndRecord(req)
		bodyRes, _ := ioutil.ReadAll(resp.Body)
		result := string(bodyRes)
		fmt.Println(result)
		if resp.StatusCode != http.StatusOK {
			t.Error("Get file fail")
		}
	})
}

func TestECUploadAndDownloadMultiCloud(t *testing.T) {
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`
{
  "TaskType": "Upload",
   "UserID": "tester",
   "DestinationPath": "path/to/jcspantest.txt",
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-hangzhou"
         },
         {
            "CloudID": "txyun-chengdu"
         }
      ],
      "N": 3,
      "K": 2
   }
}`)
		resp := sendDataAndRecord("POST", "/task", bytes.NewBuffer(jsonStr))
		var reply RequestTaskReply
		err := json.NewDecoder(resp.Body).Decode(&reply)

		token := reply.Data.Result

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()
		req, _ := postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", token)
		setCookie(req)
		sendRequestAndRecord(req)
		waitProcessorAllDone()
	})
	t.Run("Create EC Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "UserID": "tester",
     "Sid": "tttteeeesssstttt",
     "SourcePath":"path/to/jcspantest.txt",
     "SourceStoragePlan":{
        "StorageMode": "EC",
        "Clouds": [
			 {
				"CloudID": "aliyun-beijing"
			 },
			 {
				"CloudID": "aliyun-hangzhou"
			 },
			 {
				"CloudID": "txyun-chengdu"
			 }
        ],
        "N": 3,
        "K": 2
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		tid := reply.Data.Result
		logrus.Debugf("tid: %v", tid)
		waitProcessorAllDone()
	})
	var url string
	t.Run("Check File DB and get download url", func(t *testing.T) {
		url = getDownloadUrl("tester/path/to/jcspantest.txt")
		t.Logf("download url: %v", url)
		waitProcessorAllDone()
	})
	t.Run("Get file", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		resp := sendRequestAndRecord(req)
		if resp.StatusCode != http.StatusOK {
			t.Error("Get file fail")
		}
		bodyRes, _ := ioutil.ReadAll(resp.Body)
		result := string(bodyRes)
		fmt.Println(result)
	})
}
func TestReplicaUploadAndDownload(t *testing.T) {
	t.Run("Create Replica Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`{
   "TaskType": "Upload",
   "UserID": "tester",
   "DestinationPath":"path/to/jcspantest.txt",
   "DestinationStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ]
   }
}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)

		tid := reply.Data.Result

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()

		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", tid)
		setCookie(req)
		sendRequestAndRecord(req)
		waitProcessorAllDone()
	})
	t.Run("Create Replica Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "UserID": "tester",
     "SourcePath":"path/to/jcspantest.txt",
     "SourceStoragePlan":{
        "StorageMode": "Replica",
        "Clouds": [
           {
              "CloudID": "aliyun-beijing"
           },
           {
              "CloudID": "aliyun-beijing"
           },
           {
              "CloudID": "aliyun-beijing"
           }
        ]
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		downloadUrl := reply.Data.Result
		t.Logf("tid: %v", downloadUrl)
		waitProcessorAllDone()
	})
}

func TestEC2ReplicaSync(t *testing.T) {
	dstPath := "tmp/test/sync/test.jpeg"
	testECUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testECUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Sync",
   "UserID": "tester",
   "DestinationPath":"tmp/test/sync/",
   "SourcePath": "tmp/test/sync/",
   "SourceStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ],
      "N": 3,
      "K": 2
   },
   "DestinationStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ]
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	resp := sendRequestAndRecord(req)
	var reply RequestTaskReply
	_ = json.NewDecoder(resp.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitProcessorAllDone()
}

func TestReplicaMigrate(t *testing.T) {
	dstPath := "tmp/test/Migrate/test.jpeg"
	testECUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/migrate/test.txt"
	testECUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Migrate",
   "UserID": "tester",
   "DestinationPath":"",
   "SourcePath": "",
   "SourceStoragePlan":{
      "StorageMode": "Migrate",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         }
      ]
   },
   "DestinationStoragePlan":{
      "StorageMode": "Migrate",
      "Clouds": [
         {
            "CloudID": "ksyun-beijing"
         }
      ]
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	resp := sendRequestAndRecord(req)
	var reply RequestTaskReply
	_ = json.NewDecoder(resp.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitProcessorAllDone()
}

func TestReplica2ECSync(t *testing.T) {
	dstPath := "tmp/test/sync/test.jpeg"
	testReplicaUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testReplicaUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Sync",
   "UserID": "tester",
   "DestinationPath":"tmp/test/sync/",
   "SourcePath": "tmp/test/sync/",
   "SourceStoragePlan":{
      "StorageMode": "Replica",
     "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ] 
   },
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ],
      "N": 3,
      "K": 2
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	resp := sendRequestAndRecord(req)
	var reply RequestTaskReply
	_ = json.NewDecoder(resp.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitProcessorAllDone()
}

func TestReplicaUploadAndDelete(t *testing.T) {
	dstPath := "tmp/test/sync/fail_if_not_delete.jpeg"
	testReplicaUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testReplicaUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Delete",
   "UserID": "tester",
   "SourcePath": "tmp/test/sync/fail_if_not_delete.jpeg",
   "SourceStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ]
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	resp := sendRequestAndRecord(req)
	var reply RequestTaskReply
	_ = json.NewDecoder(resp.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("task fail :%v", reply.Msg)
	}
	waitProcessorAllDone()
}

func TestECUploadAndDelete(t *testing.T) {
	dstPath := "tmp/test/del/test.jpeg"
	testECUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/del/test.txt"
	testECUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Delete",
   "UserID": "tester",
   "SourcePath": "tmp/test/del/test.jpeg",
   "SourceStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         },
         {
            "CloudID": "aliyun-beijing"
         }
      ],
      "N": 3,
      "K": 2
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	resp := sendRequestAndRecord(req)
	var reply RequestTaskReply
	_ = json.NewDecoder(resp.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("task fail :%v", reply.Msg)
	}
	waitProcessorAllDone()
}

func TestMultiUpload(t *testing.T) {
	dstPath := "tmp/test/upload/test.jpeg"
	testReplicaUpload(t, dstPath, "../test/tmp/test.jpeg", "aliyun-beijing")
	dstPath = "tmp/test/upload/test.txt"
	testReplicaUpload(t, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
}

func postFile(filename string, filepath string, target_url string, token string) (*http.Request, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	writer, _ := body_writer.CreateFormField("token")
	writer.Write([]byte(token))
	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	return req, nil
}

func waitUntilAllDone(processor *TaskProcessor) {
	for true {
		time.Sleep(time.Millisecond * 500)
		if processor.taskStorage.IsAllDone() {
			return
		}
	}
}

func testECUpload(t *testing.T, dstPath string, localPath string, cloud string) {
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := fmt.Sprintf(`
{
  "TaskType": "Upload",
   "UserID": "tester",
   "DestinationPath":"%v",
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "CloudID": "%v"
         },
         {
            "CloudID": "%v"
         },
         {
            "CloudID": "%v"
         }
      ],
      "N": 3,
      "K": 2
   }
}
`, dstPath, cloud, cloud, cloud)
		jsonByte := []byte(jsonStr)

		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonByte))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		if reply.Code != http.StatusOK {
			t.Fatalf("task fail :%v", reply.Msg)
		}
		tid := reply.Data.Result
		url := fmt.Sprintf("/upload/%v", dstPath)

		req, _ = postFile("test.txt", localPath, url, tid)
		sendRequestAndRecord(req)
		waitProcessorAllDone()
	})
}

func testReplicaUpload(t *testing.T, dstPath string, localPath string, cloud string) {
	t.Run("Create Replica Upload Task and process", func(t *testing.T) {
		jsonStr := fmt.Sprintf(`
{
  "TaskType": "Upload",
   "UserID": "tester",
   "DestinationPath":"%v",
   "DestinationStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "CloudID": "%v"
         },
         {
            "CloudID": "%v"
         },
         {
            "CloudID": "%v"
         }
      ]
   }
}
`, dstPath, cloud, cloud, cloud)
		jsonByte := []byte(jsonStr)

		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonByte))
		resp := sendRequestAndRecord(req)
		var reply RequestTaskReply
		_ = json.NewDecoder(resp.Body).Decode(&reply)
		if reply.Code != http.StatusOK {
			t.Fatalf("task fail :%v", reply.Msg)
		}

		tid := reply.Data.Result
		url := fmt.Sprintf("/upload/%v", dstPath)

		req, _ = postFile("test.txt", localPath, url, tid)
		setCookie(req)
		sendRequestAndRecord(req)
		waitProcessorAllDone()
	})
}
