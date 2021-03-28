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
	"time"
)

func TestFileDatabase(t *testing.T) {
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
	collection := client.Database("Cloud").Collection("FileDatabase")
	_, err = collection.InsertOne(context.TODO(), bson.D{
		{"Id", "《1》-Wanggj"},
		{"Filename", "《1》"},
		{"Owner", "wanggj"},
		{"Size", 2},
		{"LastChange", time.Now()},
		{"SyncStatus", "oss-cn-beijing.aliyuncs.com"},
		{"ReconstructStatus", "LTAI4G3PCfrg7aXQ6EvuDo25"},
		{"DownloadUrl", "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0"},
		{"ReconstructTime", time.Now()},
	})

	mongo, err := model.NewMongoFileDatabase()
	if err != nil {
		t.Error(err)
	}

	//t.Run("test CreateFileInfo ",func(t *testing.T){
	//	file:=model.File{
	//		Id:                "《2》-Wanggj",
	//		Filename:          "《2》",
	//		Owner:             "wanggj",
	//		Size:              2,
	//		LastChange:        time.Now(),
	//		SyncStatus:        "asdfasd",
	//		ReconstructStatus: "adfasd",
	//		DownloadUrl:       "asdfasd",
	//		ReconstructTime:   time.Now(),
	//	}
	//	err:=mongo.CreateFileInfo(&file)
	//	if err!=nil{
	//		t.Error(err)
	//	}
	//})

	t.Run("test Update", func(t *testing.T) {
		file := model.File{
			Id:                "《2》-Wanggj",
			Filename:          "《2》",
			Owner:             "wanggj",
			Size:              2,
			LastChange:        time.Now(),
			SyncStatus:        "this is a new status",
			ReconstructStatus: "adfasd",
			DownloadUrl:       "asdfasd",
			ReconstructTime:   time.Now(),
		}

		err := mongo.UpdateFileInfo(&file)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("test GetInfo", func(t *testing.T) {
		file, err := mongo.GetFileInfo("tester/path/to/jcspantest.txt")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(file)
		//fmt.Println(*file)
	})

	t.Run("test Index", func(t *testing.T) {
		files, err := mongo.Index("《")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(len(files), files)
		for _, file := range files {
			fmt.Println(file)
		}
	})

	t.Run("test deleteInfo", func(t *testing.T) {
		file := model.File{
			Id:                "《2》-Wanggj",
			Filename:          "《2》",
			Owner:             "wanggj",
			Size:              2,
			LastChange:        time.Now(),
			SyncStatus:        "this is a new status",
			ReconstructStatus: "adfasd",
			DownloadUrl:       "asdfasd",
			ReconstructTime:   time.Now(),
		}
		err := mongo.DeleteFileInfo(&file)
		if err != nil {
			t.Error(err)
		}
	})

}
