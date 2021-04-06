package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type User struct {
	UserId            string           `bson:"user_id"`
	Email             string           `bson:"email"`
	Password          string           `bson:"password"`
	Nickname          string           `bson:"nickname"`
	Role              string           `bson:"role"`
	Avatar            string           `bson:"avatar"`
	LastModified      time.Time        `bson:"last_modified"`
	Preference        Preference       `bson:"preference"`
	StoragePlan       StoragePlan      `bson:"storage_plan"`
	DataStats         DataStats        `bson:"data_stats"`
	AccessCredentials AccessCredential `bson:"access_credentials"`
	Status            string           `bson:"status"`
}

type Preference struct {
	Vendor       int            `bson:"vendor"`
	StoragePrice float64        `bson:"storage_price"`
	TrafficPrice float64        `bson:"traffic_price"`
	Availability float64        `bson:"availability"`
	Latency      map[string]int `bson:"latency"`
}

type DataStats struct {
	Volume          int64            `bson:"volume"`
	UploadTraffic   map[string]int64 `bson:"upload_traffic"`
	DownloadTraffic map[string]int64 `bson:"download_traffic"`
}

type AccessCredential struct {
	CloudID  string `bson:"cloud_id"`
	UserID   string `bson:"user_id"`
	Password string `bson:"password"`
}

type UserDatabase interface {
	GetUserFromID(uid string) (*User, error)
	UpdateUserInfo(user *User) error
}

type MongoUserDatabase struct {
	databaseName   string
	clientOptions  *options.ClientOptions
	client         *mongo.Client
	collectionName string
}

func NewMongoUserDatabase() (*MongoUserDatabase, error) {
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
	return &MongoUserDatabase{
		databaseName:   util.Config.Database.DatabaseName,
		clientOptions:  clientOptions,
		client:         client,
		collectionName: "Task",
	}, nil
}

func (m *MongoUserDatabase)UpdateUserInfo(user *User) error{
	err := CheckClient(m.client, m.clientOptions)
	if err != nil {
		return err
	}

	filter:=bson.M{
		"user_id":user.UserId,
	}
	update := bson.D{
		{"$set", *user},
	}
	collection := m.client.Database(m.databaseName).Collection(m.collectionName)
	_,err=collection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		return err
	}else{
		return nil
	}
}

func (m *MongoUserDatabase)GetUserFromID(uid string) (*User, error){
	err := CheckClient(m.client, m.clientOptions)
	if err != nil {
		return nil, err
	}

	var result User

	//get the collection and find by _id
	collection := m.client.Database(m.databaseName).Collection(m.collectionName)
	err = collection.FindOne(context.TODO(), bson.D{{"user_id", uid}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &result, err
}

type InMemoryUserDatabase struct {
	UserInfoMap map[string]User
}

func NewInMemoryUserDatabase() *InMemoryUserDatabase {
	return &InMemoryUserDatabase{UserInfoMap: map[string]User{
		"tester": {
			UserId:       "tester",
			Email:        "",
			Password:     "",
			Nickname:     "",
			Role:         "",
			Avatar:       "",
			LastModified: time.Time{},
			Preference:   Preference{},
			StoragePlan:  StoragePlan{},
			DataStats: DataStats{
				Volume: 0,
				UploadTraffic: map[string]int64{
					util.Config.LocalCloudID: 0,
				},
				DownloadTraffic: map[string]int64{
					util.Config.LocalCloudID: 0,
				},
			},
			AccessCredentials: AccessCredential{},
			Status:            "",
		},
	}}
}

func (db *InMemoryUserDatabase) GetUserFromID(uid string) (*User, error) {
	user := db.UserInfoMap[uid]
	return &user, nil
}
func (db *InMemoryUserDatabase) UpdateUserInfo(user *User) error {
	db.UserInfoMap[user.UserId] = *user
	return nil
}
