package model

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type AWSBucketStorageClient struct {
	awsClient  *s3.S3
	bucketName string
}

func (client *AWSBucketStorageClient) realRemotePath(remotePath string, uid string) string {
	if remotePath[0] == '/' {
		return uid + remotePath
	} else {
		return uid + "/" + remotePath
	}
}

func GetAWSClient(endpoint, accessKeyID, secretAccessKey string) (*s3.S3, error) {
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	config := &aws.Config{
		Region:      aws.String("BEIJING"),
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(true),
		Credentials: creds,
	}
	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	client := s3.New(s)
	return client, nil
}

func (client *AWSBucketStorageClient) Upload(localPath string, remotePath string, uid string) error {
	remotePath = client.realRemotePath(remotePath, uid)
	// open file
	file, err := os.Open(localPath)
	if err != nil {
		log.WithError(err).Errorf("Open %s failed.", localPath)
		return err
	}
	defer file.Close()

	_, err = client.awsClient.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(client.bucketName),
		Key:    aws.String(remotePath),
		Body:   file,
	})

	if err != nil {
		return err
	}
	return nil
}

func (client *AWSBucketStorageClient) Download(remotePath string, localPath string, uid string) (err error) {
	remotePath = client.realRemotePath(remotePath, uid)
	file, err := os.Create(localPath)
	if err != nil {
		log.WithError(err).Errorf("os.Create(%s) failed.", localPath)
		return err
	}
	defer file.Close()
	downloader := s3manager.NewDownloaderWithClient(client.awsClient)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(client.bucketName),
		Key:    aws.String(remotePath),
	})
	return err
}

func (client *AWSBucketStorageClient) Remove(remotePath string, uid string) (err error) {
	remotePath = client.realRemotePath(remotePath, uid)
	_, err = client.awsClient.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(client.bucketName),
		Key:    aws.String(remotePath),
	})
	return err
}

func (client *AWSBucketStorageClient) Copy(srcPath string, dstPath string, uid string) (err error) {
	return errors.New("AWS Copy not yet implemented")
}

func (client *AWSBucketStorageClient) Index(remotePath string, uid string) <-chan ObjectInfo {

	return nil
}

// 递归查找
func (client *AWSBucketStorageClient) RecursiveIndex(remotePath string, uid string) <-chan ObjectInfo {

	return nil
}
func (client *AWSBucketStorageClient) GetTmpDownloadUrl(remotePath string, uid string, validTime time.Duration) (url string, err error) {

	remotePath = client.realRemotePath(remotePath, uid)
	req, _ := client.awsClient.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(client.bucketName),
		Key:    aws.String(remotePath),
	})
	urlStr, err := req.Presign(validTime)
	if err != nil {
		log.Errorf("Failed to sign request %v", err)
		return "", err
	}
	return urlStr, nil
}
