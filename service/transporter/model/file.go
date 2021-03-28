package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
	"time"
)

type File struct {
	Id                string
	Filename          string
	Owner             string
	Size              int64
	LastChange        time.Time
	SyncStatus        string // 同步状态 Pending/Deleting/Done
	ReconstructStatus string // 重建状态
	DownloadUrl       string
	ReconstructTime   time.Time
}

type FileDatabase interface {
	CreateFileInfo(file *File) (err error)
	DeleteFileInfo(file *File) (err error)
	UpdateFileInfo(file *File) (err error)
	GetFileInfo(Id string) (file *File, err error)
	Index(prefix string) (files []*File, err error)
}

type MongoFileDatabase struct {
	collectionName string
	databaseName string
	clientOption *options.ClientOptions
	client *mongo.Client
}

func NewMongoFileDatabase() (*MongoFileDatabase, error) {
	var clientOptions *options.ClientOptions
	if util.CONFIG.Database.Username!=""{
		clientOptions = options.Client().ApplyURI("mongodb://" +util.CONFIG.Database.Username+":"+util.CONFIG.Database.Password+"@"+
			util.CONFIG.Database.Host + ":" + util.CONFIG.Database.Port)
	}else{
		clientOptions = options.Client().ApplyURI("mongodb://" + util.CONFIG.Database.Host + ":" + util.CONFIG.Database.Port)
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &MongoFileDatabase{
		collectionName: "File",
		databaseName: util.CONFIG.Database.DatabaseName,
		clientOption: clientOptions,
		client: client,
	}, nil
}
func (mf *MongoFileDatabase) CreateFileInfo(file *File) (err error) {
	//check the connection
	err=CheckClient(mf.client,mf.clientOption)
	if err!=nil{
		return err
	}

	//insert the file
	collection := mf.client.Database(mf.databaseName).Collection(mf.collectionName)
	_, err = collection.InsertOne(context.TODO(), *file)
	if err != nil {
		return err
	}
	return nil
}
func (mf *MongoFileDatabase) DeleteFileInfo(file *File) (err error) {
	err=CheckClient(mf.client,mf.clientOption)
	if err!=nil{
		return err
	}

	//delete the file
	filter := bson.M{
		"id": file.Id,
	}
	collection := mf.client.Database(mf.databaseName).Collection(mf.collectionName)
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
func (mf *MongoFileDatabase) UpdateFileInfo(file *File) (err error) {
	err=CheckClient(mf.client,mf.clientOption)
	if err!=nil{
		return err
	}

	//delete the file
	filter := bson.M{
		"id": file.Id,
	}
	update := bson.D{
		{"$set", *file},
	}
	collection := mf.client.Database(mf.databaseName).Collection(mf.collectionName)
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (mf *MongoFileDatabase) GetFileInfo(Id string) (file *File, err error) {
	err=CheckClient(mf.client,mf.clientOption)
	if err!=nil{
		return nil,err
	}

	//delete the file
	var result File
	filter := bson.D{
		{"id", Id},
	}
	collection := mf.client.Database(mf.databaseName).Collection(mf.collectionName)
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (mf *MongoFileDatabase) Index(prefix string) (files []*File, err error) {
	err=CheckClient(mf.client,mf.clientOption)
	if err!=nil{
		return nil,err
	}
	//delete the file
	var result []*File
	filter := bson.M{
		"id": bson.M{
			"$regex": prefix + "*",
		},
	}
	collection := mf.client.Database(mf.databaseName).Collection(mf.collectionName)
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	//decode the result
	for cur.Next(context.TODO()) {
		var elem File
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		result = append(result, &elem)
	}
	cur.Close(context.TODO())
	return result, nil
}

type InMemoryFileDatabase struct {
	db map[string]File
}

func NewInMemoryFileDatabase() *InMemoryFileDatabase {
	return &InMemoryFileDatabase{db: make(map[string]File)}
}

func NewFileInfoFromPath(path string, uid string, fileName string) (file *File, err error) {
	fi, err := os.Stat(path)
	if fileName[0] != '/' {
		fileName = "/" + fileName
	}
	if err != nil {
		return nil, err
	}
	return &File{
		Id:                uid + fileName,
		Filename:          fileName,
		Owner:             uid,
		Size:              fi.Size(),
		LastChange:        time.Now(),
		SyncStatus:        "",
		ReconstructStatus: "",
		DownloadUrl:       "",
		ReconstructTime:   time.Time{},
	}, nil
}

func (fd *InMemoryFileDatabase) CreateFileInfo(file *File) (err error) {
	fd.db[file.Id] = *file
	return nil
}

func (fd *InMemoryFileDatabase) DeleteFileInfo(file *File) (err error) {
	delete(fd.db, file.Id)
	return nil
}

func (fd *InMemoryFileDatabase) UpdateFileInfo(file *File) (err error) {
	if _, ok := fd.db[file.Id]; !ok {
		return errors.New("file info not exist")
	}
	fd.db[file.Id] = *file
	return nil
}

func (fd *InMemoryFileDatabase) GetFileInfo(Id string) (file *File, err error) {
	f, ok := fd.db[Id]
	if ok {
		return &f, nil
	}
	return nil, errors.New("file info not exist")
}

func (fd *InMemoryFileDatabase) Index(prefix string) (files []*File, err error) {
	for _, v := range fd.db {
		if strings.HasPrefix(v.Id, prefix) {
			file := File{}
			copier.Copy(&file, v)
			files = append(files, &file)
		}
	}
	return files, nil
}

func FromFileInfoGetUidAndPath(file *File) (uid string, path string) {
	p := strings.Index(file.Id, "/")
	uid = file.Id[0:p]
	path = file.Id[p+1:]
	return
}
