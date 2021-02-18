package transporter

import (
	"github.com/minio/minio-go/v7"
	"log"
)

type S3Client struct {
	endpoint    string
	ak          string
	minioClient *minio.Client
}

type StorageDatabase interface {
	GetStorageClient(sid string, path string) StorageClient
}

type SimpleInMemoryStorageDatabase struct {
	s3ClientMap map[string]S3Client
}

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
