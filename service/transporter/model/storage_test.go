package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"testing"
	"time"
)

func TestAllStorageClient(t *testing.T) {
	err := util.ReadConfigFromFile("../transporter_config.json")
	if err != nil {
		t.Fatal("read config file fail")
	}
	clientDatabase, _ := NewMongoCloudDatabase()
	cloudsID := []string{
		"bdyun-shanghai",
		"aliyun-hangzhou",
		"ksyun-beijing",
		"txyun-chengdu",
		"txyun-guangzhou",
	}
	for _, cloudID := range cloudsID {
		client, err := clientDatabase.GetStorageClientFromName("tester", cloudID)
		if err != nil {
			t.Errorf("Cloud:%v not exist", cloudID)
			continue
		}
		err = client.Upload("../test/test.txt", "test.txt", "dev")
		if err != nil {
			t.Errorf("Cloud:%v upload error: %v", cloudID, err)
			continue
		}
		url, err := client.GetTmpDownloadUrl("test.txt", "dev", time.Hour)
		if err != nil {
			t.Errorf("Cloud:%v get tmp download url error: %v", cloudID, err)
			continue
		}
		t.Logf("%v download url: %v", cloudID, url)
	}
}

func TestIndex(t *testing.T) {
	err := util.ReadConfigFromFile("../transporter_config.json")
	if err != nil {
		t.Fatal("read config file fail")
	}
	clientDatabase, _ := NewMongoCloudDatabase()
	cloudID := "aliyun-beijing"
	client, err := clientDatabase.GetStorageClientFromName("tester", cloudID)
	if err != nil {
		t.Errorf("Cloud:%v not exist", cloudID)
	}
	objectsChan := client.Index("", "tester")
	for object := range objectsChan {
		t.Log(object.Key)
	}
}
