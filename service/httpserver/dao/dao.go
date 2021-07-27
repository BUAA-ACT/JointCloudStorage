package dao

import (
	"cloud-storage-httpserver/args"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

var (
	UserDao            *Dao
	AccessTokenDao     *Dao
	FileDao            *Dao
	TaskDao            *Dao
	VerifyCodeDao      *Dao
	MigrationAdviceDao *Dao
	CloudDao           *Dao
	VoteCloudDao       *Dao
	TempCloudDao       *Dao
	AccessKeyDao       *Dao
	DataColMap         = make(map[string]map[string]*Dao)
)

type Dao struct {
	client               *mongo.Client
	databaseConnection   *mongo.Database
	collectionConnection *mongo.Collection
	database             string
	collection           string
}

func ConnectDao(mongoURI string, database string, collection string) (*Dao, error) {
	dao := &Dao{
		database:   database,
		collection: collection,
	}
	testMap := DataColMap[database]
	if testMap == nil {
		var newCollectionMap = make(map[string]*Dao)
		DataColMap[database] = newCollectionMap
	}
	DataColMap[database][collection] = dao

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	dao.client = client
	dao.databaseConnection = client.Database(database)
	dao.collectionConnection = client.Database(database).Collection(collection)
	return dao, nil
}

func DisconnectDao(d *Dao) {
	err := d.client.Disconnect(context.TODO())
	if err != nil {
		log.Println("disconnect failed")
		log.Println(err)
	}
}

func ConnectInitDao() {
	var err error
	UserDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.UserCollection)
	if err != nil {
		log.Println(err)
	}
	AccessTokenDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.AccessTokenCollection)
	if err != nil {
		log.Println(err)
	}
	FileDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.FileCollection)
	if err != nil {
		log.Println(err)
	}
	TaskDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.TaskCollection)
	if err != nil {
		log.Println(err)
	}
	VerifyCodeDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.VerifyCodeCollection)
	if err != nil {
		log.Println(err)
	}
	MigrationAdviceDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.MigrationAdviceCollection)
	if err != nil {
		log.Println(err)
	}
	CloudDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.CloudCollection)
	if err != nil {
		log.Println(err)
	}
	VoteCloudDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.VoteCloudCollection)
	if err != nil {
		log.Println(err)
	}
	TempCloudDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.TempCloudCollection)
	if err != nil {
		log.Println(err)
	}
	AccessKeyDao, err = ConnectDao(*args.MongoURL, *args.DataBase, *args.AccessKeyCollection)
	if err != nil {
		log.Println(err)
	}
}

func AddIndex() {
	modifyAccessTokenExpireIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{"access_token_modify_time", bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(args.AccessTokenTimeOut.Seconds())),
	}
	createAccessTokenExpireIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{"access_token_create_time", bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(args.AccessTokenDateOut.Seconds())),
	}
	_, modifyAccessTokenExpireIndexErr := AccessTokenDao.collectionConnection.Indexes().CreateOne(context.TODO(), modifyAccessTokenExpireIndex)
	if modifyAccessTokenExpireIndexErr != nil {
		log.Println(modifyAccessTokenExpireIndexErr)
	}
	_, createAccessTokenExpireIndexErr := AccessTokenDao.collectionConnection.Indexes().CreateOne(context.TODO(), createAccessTokenExpireIndex)
	if createAccessTokenExpireIndexErr != nil {
		log.Println(createAccessTokenExpireIndexErr)
	}
	createVerifyCodeExpireIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{"verify_code_create_time", bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(args.AccessTokenDateOut.Seconds())),
	}
	_, createVerifyCodeExpireIndexErr := VerifyCodeDao.collectionConnection.Indexes().CreateOne(context.TODO(), createVerifyCodeExpireIndex)
	if createVerifyCodeExpireIndexErr != nil {
		log.Println(createVerifyCodeExpireIndexErr)
	}
}
