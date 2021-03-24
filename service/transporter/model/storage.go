package model

import (
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	url2 "net/url"
	"path"
	"time"
)

// 存储客户端接口，提供基础的对象存储文件管理接口
type StorageClient interface {
	Upload(localPath string, remotePath string, uid string) (err error)
	Download(remotePath string, localPath string, uid string) (err error)
	Remove(remotePath string, uid string) (err error)
	Copy(srcPath string, dstPath string, uid string) (err error)
	Index(remotePath string, uid string) <-chan ObjectInfo
	// 递归查找
	RecursiveIndex(remotePath string, uid string) <-chan ObjectInfo
	GetTmpDownloadUrl(remotePath string, uid string, validTime time.Duration) (url string, err error)
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

func (client *S3BucketStorageClient) realRemotePath(remotePath string, uid string) string {
	if remotePath[0] == '/' {
		return uid + remotePath
	} else {
		return uid + "/" + remotePath
	}
}

func (client *S3BucketStorageClient) Upload(localPath string, remotePath string, uid string) (err error) {
	ctx := context.Background()
	mime, err := mimetype.DetectFile(localPath)
	remotePath = client.realRemotePath(remotePath, uid)
	_, err = client.minioClient.FPutObject(
		ctx,
		client.bucketName,
		remotePath,
		localPath,
		minio.PutObjectOptions{ContentType: mime.String()},
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *S3BucketStorageClient) Download(remotePath string, localPath string, uid string) (err error) {
	ctx := context.Background()
	remotePath = client.realRemotePath(remotePath, uid)
	err = client.minioClient.FGetObject(ctx, client.bucketName, remotePath, localPath, minio.GetObjectOptions{})
	return err
}
func (client *S3BucketStorageClient) Remove(remotePath string, uid string) (err error) {
	ctx := context.Background()
	remotePath = client.realRemotePath(remotePath, uid)
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err = client.minioClient.RemoveObject(ctx, client.bucketName, remotePath, opts)
	return nil
}
func (client *S3BucketStorageClient) Copy(srcPath string, dstPath string, uid string) (err error) {
	ctx := context.Background()
	srcPath = client.realRemotePath(srcPath, uid)
	dstPath = client.realRemotePath(dstPath, uid)
	dst := minio.CopyDestOptions{
		Bucket: client.bucketName,
		Object: dstPath,
	}
	src := minio.CopySrcOptions{
		Bucket: client.bucketName,
		Object: srcPath,
	}
	_, err = client.minioClient.CopyObject(ctx, dst, src)
	return err
}

func (client *S3BucketStorageClient) Index(remotePath string, uid string) <-chan ObjectInfo {
	ctx := context.Background()
	remotePath = client.realRemotePath(remotePath, uid)
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

func (client *S3BucketStorageClient) RecursiveIndex(remotePath string, uid string) <-chan ObjectInfo {
	ctx := context.Background()
	remotePath = client.realRemotePath(remotePath, uid)
	ObjectCh := make(chan ObjectInfo, 1)
	go func(ObjectCh chan ObjectInfo) {
		defer close(ObjectCh)
		for obj := range client.minioClient.ListObjects(ctx, client.bucketName, minio.ListObjectsOptions{
			Prefix:    remotePath,
			Recursive: true,
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

func (client *S3BucketStorageClient) GetTmpDownloadUrl(remotePath string, uid string, validTime time.Duration) (url string, err error) {
	ctx := context.Background()
	remotePath = client.realRemotePath(remotePath, uid)
	filename := path.Base(remotePath)
	reqParams := make(url2.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+filename+"\"")
	res, err := client.minioClient.PresignedGetObject(ctx, client.bucketName, remotePath, validTime, reqParams)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
