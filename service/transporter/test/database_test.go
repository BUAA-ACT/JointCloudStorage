package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func TestDatabase(t *testing.T) {
	//insert the test data
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://192.168.105.8:20100")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	//get the collection and insert the bson
	collection:=client.Database("Cloud").Collection("Cloud")
	_,err=collection.InsertOne(context.TODO(),bson.D{
		{"id","aliyun-beijing"},
		{"storage_price",0.5},
		{"traffic_price",0.5},
		{"availability",0.9999},
		{"status","UP"},
		{"endpoint","oss-cn-beijing.aliyuncs.com"},
		{"access_key","LTAI4G3PCfrg7aXQ6EvuDo25"},
		{"secret_key","5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0"},
		{"location","116.381252,39.906569"},
	})

	mongo,err:=model.NewMongoStorageDatabase()
	if err!=nil{
		t.Error(err)
	}

	t.Run("test function ",func(t *testing.T){
		storage:=mongo.GetStorageClientFromName("asdfa","aliyun-beijing")
		fmt.Println(storage)
	})
}
