package model

import "testing"

func TestAWSBucketStorageClient(t *testing.T) {
	endpoint := "s3.fsh.bcebos.com"
	ak := "68baec42f281416ab8b6982f1a8e01be"
	sk := "2a3f515225c041ab9ccabd624afa5acc"
	bucket := "jcspan-shanghai"
	client, err := GetAWSClient(endpoint, ak, sk)
	if err != nil {
		t.Fatalf("get client err")
	}
	bucketClient := AWSBucketStorageClient{
		awsClient:  client,
		bucketName: bucket,
	}
	err = bucketClient.Upload("../test/test.txt", "test.txt", "t")
	if err != nil {
		t.Fatalf("upload err %v", err)
	}
}
