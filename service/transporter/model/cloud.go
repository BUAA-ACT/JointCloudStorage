package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type Cloud struct {
	CloudID      string  `bson:"cloud_id"`
	Endpoint     string  `bson:"endpoint"`
	AccessKey    string  `bson:"access_key" json:"-"`
	SecretKey    string  `bson:"secret_key" json:"-"`
	StoragePrice float64 `bson:"storage_price"`
	TrafficPrice float64 `bson:"traffic_price"`
	Availability float64 `bson:"availability"`
	Status       string  `bson:"status"`
	Location     string  `bson:"location"`
	Address      string  `bson:"address"`
	CloudName    string  `bson:"cloud_name"`
	ProviderName string  `bson:"provider_name"`
	Bucket       string  `bson:"bucket"`
}

// S3 客户端结构
type S3Client struct {
	name         string
	endpoint     string
	ak           string
	minioClient  *minio.Client // 已经连接好的 minio 客户端
	lastReadTime time.Time
}

// Storage 数据库
type CloudDatabase interface {
	// 通过用户的 session id 和访问路径，获取对应的 S3 客户端
	GetStorageClientFromName(uid string, name string) (StorageClient, error)
	GetCloudInfoFromCloudID(cloudID string) (*Cloud, error)
}

type MongoCloudDatabase struct {
	databaseName   string
	collectionName string
	clientOptions  *options.ClientOptions
	ClientMap      map[string]StorageClient
	ReadTimeMap    map[string]time.Time
	client         *mongo.Client
	lock           sync.Mutex
}

//get a MongoCloudDatabase
func NewMongoCloudDatabase() (*MongoCloudDatabase, error) {
	var clientOptions *options.ClientOptions
	if util.Config.Database.Username != "" {
		clientOptions = options.Client().ApplyURI("mongodb://" + util.Config.Database.Username + ":" + util.Config.Database.Password + "@" +
			util.Config.Database.Host + ":" + util.Config.Database.Port)
	} else {
		clientOptions = options.Client().ApplyURI("mongodb://" + util.Config.Database.Host + ":" + util.Config.Database.Port)
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &MongoCloudDatabase{
		databaseName:   util.Config.Database.DatabaseName,
		collectionName: "Cloud",
		clientOptions:  clientOptions,
		client:         client,
		ClientMap:      map[string]StorageClient{},
		ReadTimeMap:    map[string]time.Time{},
	}, nil
}

//update the client
func (m *MongoCloudDatabase) UpdateClient() error {
	err := m.client.Ping(context.TODO(), nil)
	if err != nil {
		clientOptions := options.Client().ApplyURI("mongodb://" + util.Config.Database.Host + ":" + util.Config.Database.Port)
		m.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

//close the client
func (m *MongoCloudDatabase) CloseClient() error {
	err := m.client.Disconnect(context.TODO())
	return err
}

func (m *MongoCloudDatabase) GetStorageClient(sid string, path string) StorageClient {
	return nil
}

func (m *MongoCloudDatabase) GetStorageClientFromName(sid string, name string) (StorageClient, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.ClientMap[name]; ok {
		if time.Now().Sub(m.ReadTimeMap[name]).Minutes() < 5 {
			if util.Config.DefaultStorageClient == util.MinioClient {
				return m.ClientMap[name].(*S3BucketStorageClient), nil
			} else {
				return m.ClientMap[name].(*AWSBucketStorageClient), nil
			}
		}
	}

	//check the client connection
	err := CheckClient(m.client, m.clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result interface{}

	//get the collection and find by _id
	collection := m.client.Database(util.Config.Database.DatabaseName).Collection("Cloud")
	err = collection.FindOne(context.TODO(), bson.D{{"cloud_id", name}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	res := result.(primitive.D).Map()
	//fmt.Println(res)
	if res != nil {
		//fmt.Println(res["endpoint"].(string))
		endpoint := res["endpoint"].(string)
		accessKeyId := res["access_key"].(string)
		secretAccessKey := res["secret_key"].(string)
		bucketName := res["bucket"].(string)
		if util.Config.DefaultStorageClient == util.MinioClient {
			minioClient, err := GetMinioClient(endpoint, accessKeyId, secretAccessKey)
			if err != nil {
				return nil, err
			}
			newClient := &S3BucketStorageClient{
				minioClient: minioClient,
				bucketName:  bucketName,
			}

			return newClient, nil
		} else {
			awsClient, err := GetAWSClient(endpoint, accessKeyId, secretAccessKey)
			if err != nil {
				log.Printf("get minio client fail: %v", err)
				return nil, err
			}
			newClient := &AWSBucketStorageClient{
				awsClient:  awsClient,
				bucketName: bucketName,
			}
			m.ClientMap["name"] = newClient
			m.ReadTimeMap["name"] = time.Now()
			return newClient, nil
		}
	} else {
		return nil, err
	}
}

func (m *MongoCloudDatabase) GetCloudInfoFromCloudID(cloudID string) (*Cloud, error) {
	//check the client connection
	err := CheckClient(m.client, m.clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result Cloud

	//get the collection and find by _id
	collection := m.client.Database(util.Config.Database.DatabaseName).Collection("Cloud")
	err = collection.FindOne(context.TODO(), bson.D{{"cloud_id", cloudID}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &result, nil
}

// 一个简单的内存 Storage 数据库
type SimpleInMemoryStorageDatabase struct {
	s3ClientMap  map[string]S3Client
	awsClientMap map[string]AWSBucketStorageClient
}

// 构造内存 Storage 数据库
func NewSimpleInMemoryStorageDatabase() *SimpleInMemoryStorageDatabase {
	endpoint := "oss-cn-beijing.aliyuncs.com"
	accessKeyID := "LTAI4G3PCfrg7aXQ6EvuDo25"
	secretAccessKey := "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0"
	bucketName := "jcspan-aliyun-bj-test"
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
	awsClient, err := GetAWSClient(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		logrus.Errorf("get aws client fail: %v", err)
		return nil
	}
	return &SimpleInMemoryStorageDatabase{
		s3ClientMap: map[string]S3Client{
			"aliyun-beijing": s3,
		},
		awsClientMap: map[string]AWSBucketStorageClient{
			"aliyun-beijing": {
				awsClient:  awsClient,
				bucketName: bucketName,
			},
		},
	}
}

func (database *SimpleInMemoryStorageDatabase) GetStorageClientFromName(uid string, name string) (StorageClient, error) {
	bucketName := "jcspan-aliyun-bj-test"
	if util.Config.DefaultStorageClient == util.MinioClient {
		return &S3BucketStorageClient{
			minioClient: database.s3ClientMap[name].minioClient,
			bucketName:  bucketName,
		}, nil
	} else {
		return &AWSBucketStorageClient{
			awsClient:  database.awsClientMap[name].awsClient,
			bucketName: database.awsClientMap[name].bucketName,
		}, nil
	}
}

func (database *SimpleInMemoryStorageDatabase) GetCloudInfoFromCloudID(cloudID string) (*Cloud, error) {
	return nil, errors.New("not impl")
}
