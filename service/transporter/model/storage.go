package model

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	url2 "net/url"
	"time"
)

// 存储客户端接口，提供基础的对象存储文件管理接口
type StorageClient interface {
	Upload(localPath string, remotePath string) (err error)
	Download(remotePath string, localPath string) (err error)
	Remove(remotePath string) (err error)
	Index(remotePath string) <-chan ObjectInfo
	GetTmpDownloadUrl(remotePath string, validTime time.Duration) (url string, err error)
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

func (client *S3BucketStorageClient) Download(remotePath string, localPath string) (err error) {
	//todo
	return nil
}
func (client *S3BucketStorageClient) Remove(remotePath string) (err error) {
	//todo
	return nil
}

func (client *S3BucketStorageClient) Index(remotePath string) <-chan ObjectInfo {
	ctx := context.Background()
	ObjectCh := make(chan ObjectInfo, 1)
	go func(ObjectCh chan ObjectInfo) {
		defer close(ObjectCh)
		for obj := range client.minioClient.ListObjects(ctx, client.bucketName, minio.ListObjectsOptions{
			Prefix:    remotePath,
			Recursive: false,
		}) {
			ObjectCh <- ObjectInfo{
				Key:          obj.Key,
				Size:         obj.Size,
				LastModified: obj.LastModified,
				ContentType:  obj.ContentType,
			}
		}
	}(ObjectCh)
	return ObjectCh
}

func (client *S3BucketStorageClient) GetTmpDownloadUrl(remotePath string, validTime time.Duration) (url string, err error) {
	ctx := context.Background()
	reqParams := make(url2.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
	res, err := client.minioClient.PresignedGetObject(ctx, client.bucketName, remotePath, validTime, reqParams)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
