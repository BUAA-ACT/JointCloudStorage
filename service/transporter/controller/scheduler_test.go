package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"testing"
	"time"
)

func TestJcsPanScheduler(t *testing.T) {
	scheduler := NewJcsPanScheduler("aliyun-hangzhou", "http://192.168.105.13:8082")
	clouds := []*model.Cloud{
		{
			CloudID:      "aliyun-hangzhou",
			Endpoint:     "",
			AccessKey:    "",
			SecretKey:    "",
			StoragePrice: 0,
			TrafficPrice: 0,
			Availability: 0,
			Status:       "",
			Location:     "",
			Address:      "localhost:8082",
		},
		{
			CloudID:      "ksyun-beijing",
			Endpoint:     "",
			AccessKey:    "",
			SecretKey:    "",
			StoragePrice: 0,
			TrafficPrice: 0,
			Availability: 0,
			Status:       "",
			Location:     "",
			Address:      "localhost:8182",
		},
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
