package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"testing"
	"time"
)

func TestJcsPanScheduler(t *testing.T) {
	cloudDatabase, _ := model.NewMongoCloudDatabase()
	scheduler := JcsPanScheduler{
		LocalCloudID:     "aliyun-hangzhou",
		SchedulerHostUrl: "http://192.168.105.13:8082",
		ReloadCloudInfo:  true,
		CloudDatabase:    cloudDatabase,
	}
	clouds := []string{
		"aliyun-hangzhou",
		"ksyun-beijing",
	}
	file := &model.File{
		FileID:            "tester/test/metaData/file.txt",
		Filename:          "file.txt",
		Owner:             "tester",
		Size:              0,
		LastModified:      time.Time{},
		SyncStatus:        "",
		ReconstructStatus: "",
		DownloadUrl:       "",
		LastReconstructed: time.Time{},
	}
	err := scheduler.UploadFileMetadata(clouds, "tester", file)
	if err != nil {
		t.Fatalf("Upload File Metadata fail:%v", err)
	}
}
