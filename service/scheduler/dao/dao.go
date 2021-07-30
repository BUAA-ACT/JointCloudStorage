package dao

import (
	"context"
	"github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Dao encapsulates database operations.
type Dao = config.ClientConfig

var globalDao *Dao

type Database struct {
	*Dao
}

func SetRealGlobalDao(realDao *Dao) {
	globalDao = realDao
}

// NewDao constructs a data access object (Dao).
func NewDao(mongoURI string, databases map[string]config.DatabaseConfig) (*Dao, error) {

	dao := Dao{}
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
	// load client
	dao.Client = client
	// get databases
	for databaseName, database := range databases {
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

func ensureIndex(collection config.CollectionConfig, index string, unique bool) error {
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
