package dao

import (
	"cloud-storage-httpserver/args"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
)

type Dao struct {
	client     *mongo.Client
	database   string
	collection string
}

func ConnectDao(mongoURI string, database string, collection string) (*Dao, error) {
	dao := &Dao{
		database:   database,
		collection: collection,
	}

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
