package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"github.com/minio/minio-go/v7"
	"testing"
	"time"
)

func TestCreateBucket(t *testing.T) {
	err := util.ReadConfigFromFile("../transporter_config.json")
	if err != nil {
		t.Fatal("read config file fail")
	}
	clientDatabase, _ := NewMongoStorageDatabase()
	t.Run("Aliyun", func(t *testing.T) {
		client, _ := clientDatabase.GetStorageClientFromName("tester", "aliyun-hangzhou")
		s3client := client.(*S3BucketStorageClient)
		err = s3client.minioClient.MakeBucket(context.TODO(), "jcspan", minio.MakeBucketOptions{
			Region:        "oss-cn-hangzhou",
			ObjectLocking: false,
		})
		if err != nil {
			t.Fatalf("make bucket err: %v", err)
		}
	})
	t.Run("ksyun", func(t *testing.T) {
		client, _ := clientDatabase.GetStorageClientFromName("tester", "ksyun-beijing")
		s3client := client.(*S3BucketStorageClient)
		err = s3client.minioClient.MakeBucket(context.TODO(), "jcspan", minio.MakeBucketOptions{
			Region:        "BEIJING",
			ObjectLocking: false,
		})
		if err != nil {
			t.Fatalf("make bucket err: %v", err)
		}
	})
	t.Run("bdyun", func(t *testing.T) {
		client, _ := clientDatabase.GetStorageClientFromName("tester", "bdyun-shanghai")
		s3client := client.(*S3BucketStorageClient)
		err = s3client.minioClient.MakeBucket(context.TODO(), "jcspan", minio.MakeBucketOptions{
			Region:        "BEIJING",
			ObjectLocking: false,
		})
		if err != nil {
			t.Fatalf("make bucket err: %v", err)
		}
	})
	t.Run("txyun-chengdu", func(t *testing.T) {
		client, _ := clientDatabase.GetStorageClientFromName("tester", "txyun-chengdu")
		s3client := client.(*S3BucketStorageClient)
		err = s3client.minioClient.MakeBucket(context.TODO(), "jcspan", minio.MakeBucketOptions{
			Region:        "ap-chengdu",
			ObjectLocking: false,
		})
		if err != nil {
			t.Fatalf("make bucket err: %v", err)
		}
	})
}

func TestAllStorageClient(t *testing.T) {
	err := util.ReadConfigFromFile("../transporter_config.json")
	if err != nil {
		t.Fatal("read config file fail")
	}
	clientDatabase, _ := NewMongoStorageDatabase()
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
