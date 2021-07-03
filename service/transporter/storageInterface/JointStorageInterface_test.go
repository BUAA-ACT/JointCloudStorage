package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var JSI *JointStorageInterface

func TestMain(m *testing.M) {
	util.ClearAll()
	util.CheckConfig()
	storage, _ := model.NewMongoTaskStorage()
	clientDatabase, _ := model.NewMongoCloudDatabase()
	fileDatabase, _ := model.NewMongoFileDatabase()
	processor := controller.TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(clientDatabase)
	// 初始化 FileInfo 数据库
	processor.FileDatabase = fileDatabase
	// 初始化 Lock
	lock, _ := controller.NewLock(util.Config.ZookeeperHost)
	processor.Lock = lock
	processor.Lock.UnLockAll("/tester")
	// 初始化 Scheduler
	// 初始化 scheduler
	scheduler := controller.JcsPanScheduler{
		LocalCloudID:     util.Config.LocalCloudID,
		SchedulerHostUrl: util.Config.SchedulerHost,
		ReloadCloudInfo:  true,
		CloudDatabase:    clientDatabase,
	}
	processor.Scheduler = &scheduler
	// 初始化 Monitor
	userDB, _ := model.NewMongoUserDatabase()
	processor.Monitor = controller.NewTrafficMonitor(userDB)
	processor.UserDatabase = userDB
	// 初始化 tempFile
	tfs, _ := util.NewTempFileStorage(util.Config.TempFilePath, time.Hour*8)
	processor.TempFileStorage = tfs

	JSI = NewInterface(&processor)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNewInterface(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	t.Log(string(recorder.Body.Bytes()))
}

func TestJointStorageInterface_PutObject(t *testing.T) {
	bodyBuf := new(bytes.Buffer)
	fh, err := os.Open("../test/test.txt")
	if err != nil {
		t.Errorf("error opening file")
	}
	io.Copy(bodyBuf, fh)
	req, err := http.NewRequest("PUT", "/test.txt", bodyBuf)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
}

func TestJointStorageInterface_GetMethod(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test.txt", nil)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Logf("body: %s", recorder.Body.String())
}

func TestJointStorageInterface_GetObjectList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Logf("body: %s", recorder.Body.String())
}
