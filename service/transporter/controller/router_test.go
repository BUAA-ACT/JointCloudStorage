package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
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

func setCookie(req *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "sid",
		Value:   "tttteeeesssstttt",
		Expires: expire,
	}
	req.AddCookie(&cookie)
}

func TestECUploadAndDownload(t *testing.T) {
	router, processor := initRouterAndProcessor()
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`
{
  "TaskType": "Upload",
   "Uid": "tester",
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
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
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
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", token)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
	t.Run("Create EC Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
		tid := reply.Data.Result
		logrus.Debugf("tid: %v", tid)
		waitUntilAllDone(processor)
	})
	var url string
	t.Run("Check File DB and get download url", func(t *testing.T) {
		fileInfo, err := processor.FileDatabase.GetFileInfo("tester/path/to/jcspantest.txt")
		if err != nil {
			t.Fatalf("get file info err:%v", err)
		}
		url = fileInfo.DownloadUrl
		if url == "" {
			t.Fatalf("get download url err")
		}
		t.Logf("download url: %v", url)
		waitUntilAllDone(processor)
	})
	t.Run("Get file", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		fmt.Println(recorder.Body)
		if recorder.Code != http.StatusOK {
			t.Error("Get file fail")
		}
	})
}

func TestECUploadAndDownloadMultiCloud(t *testing.T) {
	router, processor := initRouterAndProcessor()
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`
{
  "TaskType": "Upload",
   "Uid": "tester",
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
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		err := json.NewDecoder(recorder.Body).Decode(&reply)

		token := reply.Data.Result

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", token)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
	t.Run("Create EC Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
		tid := reply.Data.Result
		logrus.Debugf("tid: %v", tid)
		waitUntilAllDone(processor)
	})
	var url string
	t.Run("Check File DB and get download url", func(t *testing.T) {
		fileInfo, err := processor.FileDatabase.GetFileInfo("tester/path/to/jcspantest.txt")
		if err != nil {
			t.Fatalf("get file info err:%v", err)
		}
		url = fileInfo.DownloadUrl
		if url == "" {
			t.Fatalf("get download url err")
		}
		t.Logf("download url: %v", url)
		waitUntilAllDone(processor)
	})
	t.Run("Get file", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		fmt.Println(recorder.Body)
		if recorder.Code != http.StatusOK {
			t.Error("Get file fail")
		}
	})
}
func TestReplicaUploadAndDownload(t *testing.T) {
	router, processor := initRouterAndProcessor()
	t.Run("Create Replica Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`{
   "TaskType": "Upload",
   "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)

		tid := reply.Data.Result

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()

		//req,_ http.Post("/upload/jcspan/path/to/file", "multipart/form-data", body)
		//	req, _ = http.NewRequest("POST", "/upload/jcspan/path/to/file", body)
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", tid)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
	t.Run("Create Replica Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
		downloadUrl := reply.Data.Result
		t.Logf("tid: %v", downloadUrl)
		waitUntilAllDone(processor)
	})
}

func TestEC2ReplicaSync(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/sync/未命名.png"
	testECUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testECUpload(t, router, processor, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Sync",
   "Uid": "tester",
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var reply RequestTaskReply
	_ = json.NewDecoder(recorder.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitUntilAllDone(processor)
}

func TestReplicaMigrate(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/Migrate/未命名.png"
	testECUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/migrate/test.txt"
	testECUpload(t, router, processor, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Migrate",
   "Uid": "tester",
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var reply RequestTaskReply
	_ = json.NewDecoder(recorder.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitUntilAllDone(processor)
}

func TestReplica2ECSync(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/sync/未命名.png"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Sync",
   "Uid": "tester",
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var reply RequestTaskReply
	_ = json.NewDecoder(recorder.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("Sync fail :%v", reply.Msg)
	}
	waitUntilAllDone(processor)
}

func TestReplicaUploadAndDelete(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/sync/未命名.png"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/sync/test.txt"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Delete",
   "Uid": "tester",
   "SourcePath": "tmp/test/sync/未命名.png",
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var reply RequestTaskReply
	_ = json.NewDecoder(recorder.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("task fail :%v", reply.Msg)
	}
	waitUntilAllDone(processor)
}

func TestECUploadAndDelete(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/del/未命名.png"
	testECUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/del/test.txt"
	testECUpload(t, router, processor, dstPath, "../test/tmp/test.txt", "aliyun-beijing")
	jsonStr := []byte(`
{
  "TaskType": "Delete",
   "Uid": "tester",
   "SourcePath": "tmp/test/del/未命名.png",
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var reply RequestTaskReply
	_ = json.NewDecoder(recorder.Body).Decode(&reply)
	if reply.Code != http.StatusOK {
		t.Fatalf("task fail :%v", reply.Msg)
	}
	waitUntilAllDone(processor)
}

func TestMultiUpload(t *testing.T) {
	router, processor := initRouterAndProcessor()
	dstPath := "tmp/test/upload/未命名.png"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/未命名.png", "aliyun-beijing")
	dstPath = "tmp/test/upload/未命名.png"
	testReplicaUpload(t, router, processor, dstPath, "../test/tmp/未命名1.png", "aliyun-beijing")
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

func testECUpload(t *testing.T, router *Router, processor *TaskProcessor, dstPath string, localPath string, cloud string) {
	t.Run("Create EC Upload Task and process", func(t *testing.T) {
		jsonStr := fmt.Sprintf(`
{
  "TaskType": "Upload",
   "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
		if reply.Code != http.StatusOK {
			t.Fatalf("task fail :%v", reply.Msg)
		}
		tid := reply.Data.Result
		url := fmt.Sprintf("/upload/%v", dstPath)

		req, _ = postFile("test.txt", localPath, url, tid)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
}

func testReplicaUpload(t *testing.T, router *Router, processor *TaskProcessor, dstPath string, localPath string, cloud string) {
	t.Run("Create Replica Upload Task and process", func(t *testing.T) {
		jsonStr := fmt.Sprintf(`
{
  "TaskType": "Upload",
   "Uid": "tester",
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
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		var reply RequestTaskReply
		_ = json.NewDecoder(recorder.Body).Decode(&reply)
		if reply.Code != http.StatusOK {
			t.Fatalf("task fail :%v", reply.Msg)
		}

		tid := reply.Data.Result
		url := fmt.Sprintf("/upload/%v", dstPath)

		req, _ = postFile("test.txt", localPath, url, tid)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
}
