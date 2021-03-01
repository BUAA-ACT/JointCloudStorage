package model

import (
	"github.com/minio/minio-go/v7"
	"log"
)

// S3 客户端结构
type S3Client struct {
	endpoint    string
	ak          string
	minioClient *minio.Client // 已经连接好的 minio 客户端
}

// Storage 数据库
type StorageDatabase interface {
	// 通过用户的 session id 和访问路径，获取对应的 S3 客户端
	GetStorageClient(sid string, path string) StorageClient
}

// 一个简单的内存 Storage 数据库
type SimpleInMemoryStorageDatabase struct {
	s3ClientMap map[string]S3Client
}

// 构造内存 Storage 数据库
func NewSimpleInMemoryStorageDatabase() *SimpleInMemoryStorageDatabase {
	endpoint := "oss-cn-beijing.aliyuncs.com"
	accessKeyID := "LTAI4G3PCfrg7aXQ6EvuDo25"
	secretAccessKey := "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0"

	minioClient, err := GetMinioClient(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		log.Panicf("get minio client fail: %v", err)
		return nil
	}
	s3 := S3Client{
		endpoint:    endpoint,
		ak:          accessKeyID,
		minioClient: minioClient,
	}
	return &SimpleInMemoryStorageDatabase{
		s3ClientMap: map[string]S3Client{
			endpoint: s3,
		},
	}
}

func (database *SimpleInMemoryStorageDatabase) GetStorageClient(sid string, path string) StorageClient {
	s3Name := "oss-cn-beijing.aliyuncs.com"
	bucketName := "jcspan-aliyun-bj-test"
	return &S3BucketStorageClient{
		minioClient: database.s3ClientMap[s3Name].minioClient,
		bucketName:  bucketName,
	}
}
