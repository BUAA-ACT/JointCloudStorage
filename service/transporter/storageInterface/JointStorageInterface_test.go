package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewInterface(t *testing.T) {
	util.ClearAll()
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
	//scheduler := JcsPanScheduler{
	//	LocalCloudID:     "aliyun-hohhot",
	//	SchedulerHostUrl: "http://192.168.105.13:8082",
	//	ReloadCloudInfo:  true,
	//	CloudDatabase:    clientDatabase,
	//}
	processor.Scheduler = &scheduler
	// 初始化 Monitor
	userDB, _ := model.NewMongoUserDatabase()
	processor.Monitor = controller.NewTrafficMonitor(userDB)
	processor.UserDatabase = userDB

	jsi := NewInterface(&processor)
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	jsi.ServeHTTP(recorder, req)
	t.Log(string(recorder.Body.Bytes()))
}
