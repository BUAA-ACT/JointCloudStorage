package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//type baseTask struct {
//	Id 				primitive.ObjectID			`bson:"_id,omitempty"`
//	TaskType        TaskType
//	State           TaskState
//	StartTime       time.Time
//	Uid             string
//	SourcePath      string
//	DestinationPath string
//	TaskOptions     *TaskOptions
//}
//
////change a baseTask to Task
//func bt2t(task baseTask ) *Task{
//	return &Task{
//		Tid: 			string(task.Id[:]),
//		TaskType: 		task.TaskType,
//		State:			task.State,
//		StartTime: 		task.StartTime,
//		Uid: 			task.Uid,
//		SourcePath: 	task.SourcePath,
//		DestinationPath: task.DestinationPath,
//		TaskOptions: 	task.TaskOptions,
//	}
//}
//
////change task to basetask,Id field of the baseTask is nil
//func t2bt(task Task) *baseTask{
//	return &baseTask{
//		Id:				nil,
//		TaskType: 		task.TaskType,
//		State:			task.State,
//		StartTime: 		task.StartTime,
//		Uid: 			task.Uid,
//		SourcePath: 	task.SourcePath,
//		DestinationPath: task.DestinationPath,
//		TaskOptions: 	task.TaskOptions,
//	}
//}

type MongoTaskStorage struct {
	databaseName string
	clientOptions *options.ClientOptions
	client *mongo.Client
	collectionName string
	maxTid int
}

//func NewMongoTaskStorage() *MongoTaskStorage
//create a struct MongoTaskStorage
func NewMongoTaskStorage() (*MongoTaskStorage, error) {
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
	return &MongoTaskStorage{
		databaseName: util.CONFIG.Database.DatabaseName,
		clientOptions: clientOptions,
		client: client,
		collectionName: "Tasks",
		maxTid: 0,
	}, nil
}

//func (task *MongoTaskStorage)AddTask(t *baseTask)(tid int,err error)
//insert a task into the table
func (task *MongoTaskStorage) AddTask(t *Task) (primitive.ObjectID, error) {
	if t.State != BLOCKED {
		t.State = WAITING
	}
	//check the connection
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		return primitive.NewObjectID(),err
	}

	//get the collection and insert the bson
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	insertResult, err := collection.InsertOne(context.TODO(), *t)
	objectID := insertResult.InsertedID.(primitive.ObjectID)

	return objectID, nil
}

//get at most n tasks with state WAITTING
func (task *MongoTaskStorage) GetTaskList(n int) (t []*Task) {
	//check the connection
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		log.Println(err)
		return nil
	}

	//find the table
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	findOptions := options.Find()
	findOptions.SetLimit(int64(n))
	cur, err := collection.Find(context.TODO(), bson.D{{"state", WAITING}}, findOptions)

	if err != nil {
		log.Print(err)
		return nil
	}

	//change the result
	for cur.Next(context.TODO()) {
		var elem Task
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		t = append(t, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Print(err)
		return nil
	}

	//close the cursor
	cur.Close(context.TODO())
	return
}

//get one task by _id
func (task *MongoTaskStorage) GetTask(tid primitive.ObjectID) (*Task, error) {
	//check the connection
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		return nil,err
	}

	var result Task

	//get the collection and find by _id
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	err = collection.FindOne(context.TODO(), bson.D{{"_id", tid}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &result, err
}

//set the state to task tid to TaskState
func (task *MongoTaskStorage) SetTaskState(tid primitive.ObjectID, state TaskState) error {
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		log.Println(err)
		return err
	}

	filter := bson.D{{"_id", tid}}
	update := bson.D{
		{"$set", bson.D{
			{"state", state},
		}},
	}
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

//set the task by _id
func (task *MongoTaskStorage) SetTask(tid primitive.ObjectID, t *Task) error {
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		return err
	}

	update := bson.D{
		{"$set", bson.D{
			{"taskType", t.TaskType},
			{"state", t.State},
			{"startTime", t.StartTime},
			{"uid", t.Uid},
			{"sourcePath", t.SourcePath},
			{"destinationPath", t.DestinationPath},
			{"taskOptions", t.TaskOptions},
		}},
	}
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	_, err = collection.UpdateByID(context.TODO(), tid, update)
	if err != nil {
		return err
	}
	return nil
}

//delete the task
func (task *MongoTaskStorage) DelTask(tid primitive.ObjectID) error {
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		return err
	}

	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	_, err = collection.DeleteOne(context.TODO(), bson.D{{"id", tid}})
	if err != nil {
		return nil
	}
	return nil
}


func (task *MongoTaskStorage) IsAllDone() bool {
	//check the client
	err:=CheckClient(task.client,task.clientOptions)
	if err!=nil{
		log.Println(err)
		return false
	}

	//get the collection and find by _id
	collection := task.client.Database(task.databaseName).Collection(task.collectionName)
	filter := bson.M{"state": bson.M{
		"$nin": bson.A{FAIL, FINISH},
	}}
	result, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Print(err)
		return false
	}
	for result.Next(context.TODO()) {
		return false
	}
	result.Close(context.TODO())
	return true
}
