package dao

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Dao encapsulates database operations.
type Dao = ClientConfig

type CollectionConfig struct {
	CollectionHandler *mongo.Collection
}

type DatabaseConfig struct {
	DatabaseHandler *mongo.Database
	Collections     map[string]*CollectionConfig
}

type ClientConfig struct {
	Client    *mongo.Client
	Databases map[string]*DatabaseConfig
}

type Database struct {
	*Dao
}

//func SetRealGlobalDao(realDao *Dao) {
//	globalDao = realDao
//}

var (
	Clients        map[string]*ClientConfig
	DaoClientsLock sync.RWMutex
)

func NewClient(mongoURI string) (*mongo.Client, error) {
	// construct client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	cancelFunc()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
}

// NewDao constructs a data access object (Dao).
func NewDao(mongoURI string, databases map[string]*DatabaseConfig) (*Dao, error) {
	config.GetConfig()

	DaoClientsLock.Lock()
	dao := Clients[mongoURI]
	// get this client
	if dao == nil {
		dao = &Dao{
			Client:    nil,
			Databases: map[string]*DatabaseConfig{},
		}
		// load into whole dao
		Clients[mongoURI] = dao
	}
	// complete client
	if dao.Client == nil {
		// connect with URI
		client, err := NewClient(mongoURI)
		if err != nil {
			return nil, err
		}
		dao.Client = client
	}
	// handle with databases
	if dao.Databases == nil {
		dao.Databases = map[string]*DatabaseConfig{}
	}
	// get the databases need to be established
	for registerDatabaseName, registerDatabase := range databases {
		// get the database from whole dao
		if dao.Databases[databaseName] == nil {
			// new database pointer
			database = &DatabaseConfig{
				DatabaseHandler: nil,
				Collections:     map[string]*CollectionConfig{},
			}
		}
		if database.Collections == nil {
			database.Collections = map[string]*CollectionConfig{}
		}
		// complete with the database
		if database.DatabaseHandler == nil {

		}
		// add database handler into dao
		database.DatabaseHandler = client.Database(databaseName)

		collections := database.Collections
		for collectionName, collection := range collections {
			collection.CollectionHandler = database.DatabaseHandler.Collection(collectionName)
			// TODO: ensure user and file
			err = ensureIndex(collection, "cloud_id", true)
			if err != nil {
				logrus.Println(err)
			}
			// load collection
			collections[collectionName] = collection
		}
		//load database
		databases[databaseName] = database
	}
	dao.Databases = databases
	return &dao, nil
}

func ensureIndex(collection CollectionConfig, index string, unique bool) error {
	idx := mongo.IndexModel{
		Keys: bson.M{
			index: 1,
		},
		Options: &options.IndexOptions{
			Unique: &unique,
		},
	}

	_, err := collection.CollectionHandler.Indexes().CreateOne(context.TODO(), idx)
	if err != nil {
		return err
	}
	return nil
}

func GetDatabaseInstance() Database {
	if globalDao == nil {
		conf := config.GetConfig()
		dao, err := NewDao(conf.FlagMongo, conf.FlagEnv)
		if err != nil {
			logrus.Errorf("创建 Dao 失败： %v", err)
		}
		globalDao = dao
		logrus.Infof("创建全局 Dao 实例成功")
	}
	return Database{globalDao}
}
