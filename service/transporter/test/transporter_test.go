package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter"
	"context"
	"github.com/minio/minio-go/v7"
	"log"
	"net/url"
	"testing"
	"time"
)


func TestConnection(t *testing.T) {
	ctx := context.Background()
	endpoint := "oss-cn-beijing.aliyuncs.com"
	accessKeyID := "LTAI4G3PCfrg7aXQ6EvuDo25"
	secretAccessKey := "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0"
	bucketName := "jcspan-aliyun-bj-test"

	client, err := transporter.GetMinioClient(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		t.Errorf("get MinioClient error: %v", err)
	}
	t.Run("test BucketExists", func(t *testing.T) {
		isExist, _ := client.BucketExists(ctx, bucketName)
		if !isExist {
			t.Errorf("%v not exist", bucketName)
		}
	})
	t.Run ("test update", func(t *testing.T) {
		objectName := "test.txt"
		filePath := "/tmp/test.txt"
		contentType := "text/plain"
		n, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			log.Fatalln(err)
		}
		t.Logf("Successfully uploaded %s of size %v\n", objectName, n)
	})
	t.Run("test presigned get object", func(t *testing.T) {
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
		objects := client.ListObjects(ctx, bucketName,minio.ListObjectsOptions{Recursive:true})
		object := <- objects
		t.Logf("%v", object)
		res, err := client.PresignedGetObject(ctx, bucketName, object.Key, time.Minute*5, reqParams)
		if err != nil {
			t.Errorf("Presigned Get Object fail: %v", err)
		}else{
			t.Logf("Successfully generated presigned URL: %v", res)
		}
	})

}