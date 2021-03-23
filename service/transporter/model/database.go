package model

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// S3 客户端结构
type S3Client struct {
	name         string
	endpoint     string
	ak           string
	minioClient  *minio.Client // 已经连接好的 minio 客户端
	lastReadTime time.Time
}

// Storage 数据库
type StorageDatabase interface {
	// 通过用户的 session id 和访问路径，获取对应的 S3 客户端
	GetStorageClientFromName(uid string, name string) StorageClient
}

type MongoStorageDatabase struct {
	s3ClientMap map[string]S3Client
	client      *mongo.Client
}

//get a MongoStorageDatabase
func NewMongoStorageDatabase() (*MongoStorageDatabase, error) {
	clientOptions := options.Client().ApplyURI("mongodb://192.168.105.8:20100")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &MongoStorageDatabase{
		client:      client,
		s3ClientMap: map[string]S3Client{},
	}, nil
}

//update the client
func (m *MongoStorageDatabase) UpdateClient() error {
	err := m.client.Ping(context.TODO(), nil)
	if err != nil {
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		m.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

//close the client
func (m *MongoStorageDatabase) CloseClient() error {
	err := m.client.Disconnect(context.TODO())
	return err
}

func (m *MongoStorageDatabase) GetStorageClient(sid string, path string) StorageClient {
	return nil
}

func (m *MongoStorageDatabase) GetStorageClientFromName(sid string, name string) StorageClient {
	bucketName := "jcspan-aliyun-bj-test"
	if s3Client, ok := m.s3ClientMap[name]; ok {
		if time.Now().Sub(s3Client.lastReadTime).Minutes() < 5 {
			return &S3BucketStorageClient{
				bucketName:  bucketName,
				minioClient: s3Client.minioClient,
			}
		}
	}

	//check the client connection
	err := m.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		m.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	var result interface{}

	//get the collection and find by _id
	collection := m.client.Database("Cloud").Collection("Cloud")
	err = collection.FindOne(context.TODO(), bson.D{{"id", name}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil
	}
	res := result.(primitive.D).Map()
	fmt.Println(res)
	if res != nil {
		fmt.Println(res["endpoint"].(string))
		endpoint := res["endpoint"].(string)
		accessKeyId := res["access_key"].(string)
		secretAccessKey := res["secret_key"].(string)
		minioClient, err := GetMinioClient(endpoint, accessKeyId, secretAccessKey)
		if err != nil {
			log.Panicf("get minio client fail: %v", err)
			return nil
		}
		s3 := S3Client{
			name:         "aliyun-beijing",
			endpoint:     res["endpoint"].(string),
			ak:           res["access_key"].(string),
			minioClient:  minioClient,
			lastReadTime: time.Now(),
		}
		m.s3ClientMap[name] = s3
		return &S3BucketStorageClient{
			bucketName:  bucketName,
			minioClient: s3.minioClient,
		}
	} else {
		return nil
	}
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
		name:        "aliyun-beijing",
		endpoint:    endpoint,
		ak:          accessKeyID,
		minioClient: minioClient,
	}
	return &SimpleInMemoryStorageDatabase{
		s3ClientMap: map[string]S3Client{
			"aliyun-beijing": s3,
		},
	}
}

func (database *SimpleInMemoryStorageDatabase) GetStorageClientFromName(uid string, name string) StorageClient {
	bucketName := "jcspan-aliyun-bj-test"
	return &S3BucketStorageClient{
		minioClient: database.s3ClientMap[name].minioClient,
		bucketName:  bucketName,
	}
}
