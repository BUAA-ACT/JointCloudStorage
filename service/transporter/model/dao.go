package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func CheckClient(client *mongo.Client, connectionOptions *options.ClientOptions) error {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		client, err = mongo.Connect(context.TODO(), connectionOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

type CollectionName string

const (
	AccessKeyCollection CollectionName = "AccessKey"
)

// Dao encapsulates database operations.
type Dao struct {
	Client              *mongo.Client
	database            string
	accessKeyCollection string
	fileCollection      string
	userCollection      string
	migrationAdvice     string
}

func (d *Dao) TestPing() (ok bool) {
	err := d.Client.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true
}

// NewDao constructs a data access object (Dao).
func NewDao(mongoURI, database, accessKeyCollection string) (*Dao, error) {
	dao := &Dao{
		database:            database,
		accessKeyCollection: accessKeyCollection,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	dao.Client = client
	return dao, nil
}

func InitDao() (*Dao, error) {
	mongoURL := fmt.Sprintf("mongodb://%s:%s", util.Config.Database.Host, util.Config.Database.Port) // todo 支持密码
	return NewDao(mongoURL, util.Config.Database.DatabaseName, "AccessKey")
}

// GetCollection 获取数据库 Collection 连接
// name Collection 名称
func (d *Dao) GetCollection(name CollectionName) (collection *mongo.Collection) {
	switch name {
	case AccessKeyCollection:
		return d.Client.Database(d.database).Collection(d.accessKeyCollection)
	default:
		util.Log(logrus.ErrorLevel, "dao getCollection", "get collection not exist",
			"collection name", string(name), "")
		return nil
	}
}
