package model

import (
	"context"
	"fmt"
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
	client *mongo.Client
	maxTid int
}

//func NewMongoTaskStorage() *MongoTaskStorage
//create a struct MongoTaskStorage
func NewMongoTaskStorage() (*MongoTaskStorage, error) {
	clientOptions := options.Client().ApplyURI("mongodb://192.168.105.8:20100")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &MongoTaskStorage{
		client: client,
		maxTid: 0,
	}, nil
}

//func (task *MongoTaskStorage)AddTask(t *baseTask)(tid int,err error)
//insert a task into the table
func (task *MongoTaskStorage) AddTask(t *Task) (tid primitive.ObjectID, err error) {
	//check the connection
	err = task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return primitive.NewObjectID(), err
		}
	}

	//get the collection and insert the bson
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	insertResult, err := collection.InsertOne(context.TODO(), *t)
	objectID := insertResult.InsertedID.(primitive.ObjectID)

	return objectID, nil
}

//get at most n tasks with state WAITTING
func (task *MongoTaskStorage) GetTaskList(n int) (t []*Task) {
	//check the connection
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return nil
		}
	}

	//find the table
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	findOptions := options.Find()
	findOptions.SetLimit(int64(n))
	cur, err := collection.Find(context.TODO(), bson.D{{"State", WAITING}}, findOptions)

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
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return nil, err
		}
	}

	var result Task

	//get the collection and find by _id
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	err = collection.FindOne(context.TODO(), bson.D{{"_id", tid}}).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &result, err
}

//set the state to task tid to TaskState
func (task *MongoTaskStorage) SetTaskState(tid primitive.ObjectID, state TaskState) error {
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
	}

	filter := bson.D{{"_id", tid}}
	update := bson.D{
		{"$set", bson.D{
			{"state", state},
		}},
	}
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

//set the task by _id
func (task *MongoTaskStorage) SetTask(tid primitive.ObjectID, t *Task) error {
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
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
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	_, err = collection.UpdateByID(context.TODO(), tid, update)
	if err != nil {
		return err
	}
	return nil
}

//delete the task
func (task *MongoTaskStorage) DelTask(tid primitive.ObjectID) error {
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
	}

	collection := task.client.Database("transporterTasks").Collection("Tasks")
	_, err = collection.DeleteOne(context.TODO(), bson.D{{"id", tid}})
	if err != nil {
		return nil
	}
	return nil
}

func (task *MongoTaskStorage) UpdateClient() error {
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (task *MongoTaskStorage) IsAllDone() bool {
	//check the client
	err := task.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)
		clientOptions := options.Client().ApplyURI("mongodb://192.168.106.8:20100")
		task.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Println(err)
			return false
		}
	}

	//get the collection and find by _id
	collection := task.client.Database("transporterTasks").Collection("Tasks")
	filter := bson.M{"State": bson.M{
		"$nin": bson.A{FAIL, FINISH},
	}}
	result, err := collection.Find(context.TODO(), filter)
	fmt.Println(*result)
	if err != nil {
		log.Print(err)
		return false
	}
	if result != nil {
		return false
	}
	return true
}
