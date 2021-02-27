package transporter

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// 存储客户端接口，提供基础的对象存储文件管理接口
type StorageClient interface {
	Upload(localPath string, remotePath string) (err error)
	download(remotePath string, localPath string) (err error)
	Remove(remotePath string) (err error)
}

type S3BucketStorageClient struct {
	minioClient *minio.Client
	bucketName  string
}

func GetMinioClient(endpoint, accessKeyID, secretAccessKey string) (*minio.Client, error) {
	//ctx := context.Background()
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return minioClient, err
}

func NewS3BucketStorageClient(client *minio.Client) (*S3BucketStorageClient, error) {
	return &S3BucketStorageClient{
		minioClient: client,
		bucketName:  "",
	}, nil
}

func (client *S3BucketStorageClient) Upload(localPath string, remotePath string) (err error) {
	ctx := context.Background()
	//mime, err := mimetype.DetectFile(localPath)
	_, err = client.minioClient.FPutObject(
		ctx,
		client.bucketName,
		remotePath,
		localPath,
		minio.PutObjectOptions{ContentType: "text/plain"}, //todo
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *S3BucketStorageClient) download(remotePath string, localPath string) (err error) {
	//todo
	return nil
}
func (client *S3BucketStorageClient) Remove(remotePath string) (err error) {
	//todo
	return nil
}
