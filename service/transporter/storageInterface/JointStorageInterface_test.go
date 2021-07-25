package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var JSI *JointStorageInterface
var AK, SK string

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
	// 初始化 AccessDB
	dao, _ := model.InitDao()
	accessKeyDB := model.AccessKeyDB{Dao: dao}
	processor.AccessKeyDatabase = &accessKeyDB
	processor.StartProcessTasks(context.Background())

	key, _ := processor.AccessKeyDatabase.GenerateKeys("jsitest")
	AK = key.AccessKey
	SK = key.SecretKey
	JSI = NewInterface(&processor)

	genTestFile()
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

func TestJointStorageInterface_GetServerInfo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/state/server", nil)
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Log(recorder.Body.String())
}

func TestJointStorageInterface_PutObject(t *testing.T) {
	bodyBuf := new(bytes.Buffer)
	fh, err := os.Open("../test/tmp/test.txt")
	if err != nil {
		t.Errorf("error opening file")
	}
	io.Copy(bodyBuf, fh)
	req, err := http.NewRequest("PUT", "/object/jsiTest.txt", bodyBuf)
	req, err = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
}

func TestJointStorageInterface_GetMethod(t *testing.T) {
	req, _ := http.NewRequest("GET", "/object/jsiTest.txt", nil)
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Logf("body: %s", recorder.Body.String())
}

func TestJointStorageInterface_GetObjectList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/object/", nil)
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Logf("body: %s", recorder.Body.String())
}

func TestJointStorageInterface_DeleteObject(t *testing.T) {
	TestJointStorageInterface_PutObject(t)
	req, _ := http.NewRequest("DELETE", "/object/jsiTest.txt", nil)
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
}

func TestJointStorageInterface_GetStorageInfo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/state/storage", nil)
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("http code incorrect")
	}
	t.Logf("StorageInfo: %s", recorder.Body.String())
}

func TestJointStorageInterface_PostStoragePlan(t *testing.T) {
	TestJointStorageInterface_PutObject(t)
	TestJointStorageInterface_GetMethod(t)
	waitAllDone()
	jsonStr := []byte(`
{
    "StorageMode": "Replica",
    "Clouds": [
        "aliyun-hangzhou","aliyun-hohhot","aliyun-qingdao"
    ],
    "N": 3,
    "K": 1
}`)
	req, _ := http.NewRequest("POST", "/state/plan", bytes.NewBuffer(jsonStr))
	req, _ = JSISign(req, AK, SK)
	recorder := httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	taskId := recorder.Body.String()
	t.Logf("获取到迁移 taskID: %s", taskId)
	JSI.processor.TaskStorage.IsAllDone()
	waitAllDone()
	TestJointStorageInterface_GetMethod(t)
	jsonStr = []byte(`
{
    "StorageMode": "EC",
    "Clouds": [
        "aliyun-qingdao","aliyun-hohhot","aliyun-hangzhou"
    ],
    "N": 3,
    "K": 2
}`)
	req, _ = http.NewRequest("POST", "/state/plan", bytes.NewBuffer(jsonStr))
	req, _ = JSISign(req, AK, SK)
	recorder = httptest.NewRecorder()
	JSI.ServeHTTP(recorder, req)
	taskId = recorder.Body.String()
	t.Logf("获取到迁移 taskID: %v", taskId)
	JSI.processor.TaskStorage.IsAllDone()
	waitAllDone()
	TestJointStorageInterface_GetMethod(t)
}

func waitAllDone() {
	for !JSI.processor.TaskStorage.IsAllDone() {
		time.Sleep(500 * time.Millisecond)
	}
}
