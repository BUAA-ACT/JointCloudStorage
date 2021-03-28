package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"bytes"
	"context"
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

func initRouterAndProcessor() (*Router, *TaskProcessor) {
	var storage model.TaskStorage
	var clientDatabase model.StorageDatabase
	var fileDatabase model.FileDatabase
	err := util.CheckConfig()
	if err != nil {
		return nil, nil
	}
	if util.CONFIG.Database.Driver == util.MongoDB {
		util.ClearAll()
		storage, _ = model.NewMongoTaskStorage()
		clientDatabase, _ = model.NewMongoStorageDatabase()
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
	// 初始化路由
	router := NewTestRouter(processor)
	// 启动 processor
	processor.StartProcessTasks(context.Background())
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
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
   "DestinationPath":"path/to/upload/",
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ],
      "N": 2,
      "K": 1
   }
}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		token := recorder.Body.String()

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
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           }
        ],
        "N": 2,
        "K": 1
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		tid := recorder.Body.String()
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
   "DestinationPath":"path/to/upload/",
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-hangzhou"
         },
         {
            "ID": "txyun-chengdu"
         }
      ],
      "N": 2,
      "K": 1
   }
}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		token := recorder.Body.String()

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
				"ID": "aliyun-beijing"
			 },
			 {
				"ID": "aliyun-hangzhou"
			 },
			 {
				"ID": "txyun-chengdu"
			 }
        ],
        "N": 2,
        "K": 1
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		tid := recorder.Body.String()
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
   "DestinationPath":"/path/to/upload/",
   "DestinationStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ]
   }
}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		tid := recorder.Body.String()

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
		time.Sleep(time.Second * 5)
	})
	t.Run("Create Replica Download Task", func(t *testing.T) {
		jsonStr := []byte(`
  {
    "TaskType": "Download",
     "Uid": "tester",
     "Sid": "tttteeeesssstttt",
     "SourcePath":"path/to/jcspantest.txt",
     "SourceStoragePlan":{
        "StorageMode": "Replica",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           }
        ]
     }
  }
`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		downloadUrl := recorder.Body.String()
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
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ],
      "N": 2,
      "K": 1
   },
   "DestinationStoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ]
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
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
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ] 
   },
   "DestinationStoragePlan":{
      "StorageMode": "EC",
      "Clouds": [
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ],
      "N": 2,
      "K": 1
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
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
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ]
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
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
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ],
      "N": 2,
      "K": 1
   }
}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
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
            "ID": "%v"
         },
         {
            "ID": "%v"
         },
         {
            "ID": "%v"
         }
      ],
      "N": 2,
      "K": 1
   }
}
`, dstPath, cloud, cloud, cloud)
		jsonByte := []byte(jsonStr)

		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonByte))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		tid := recorder.Body.String()
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
            "ID": "%v"
         },
         {
            "ID": "%v"
         },
         {
            "ID": "%v"
         }
      ]
   }
}
`, dstPath, cloud, cloud, cloud)
		jsonByte := []byte(jsonStr)

		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonByte))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		tid := recorder.Body.String()
		url := fmt.Sprintf("/upload/%v", dstPath)

		req, _ = postFile("test.txt", localPath, url, tid)
		setCookie(req)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		waitUntilAllDone(processor)
	})
}
