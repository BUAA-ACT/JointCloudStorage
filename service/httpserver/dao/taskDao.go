package dao

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) GetTask(taskId string, userID string, isSingle bool) (*[]model.Task, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var filter interface{}
	filterSingle := bson.M{
		"task_id": taskId,
	}
	filterAll := bson.M{
		"user_id": userID,
	}
	if isSingle {
		filter = filterSingle
	} else {
		filter = filterAll
	}
	var tasks []model.Task = make([]model.Task, 0)
	result, findErr := col.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var task model.Task
		decodeErr := result.Decode(&task)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		tasks = append(tasks, task)
	}
	return &tasks, true
}

func (d *Dao) GetUserMigrate(userID string) (*model.Task, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":   userID,
		"task_type": args.TaskTypeMigrate,
	}
	var migrationTask model.Task
	decodeErr := col.FindOne(context.TODO(), filter).Decode(&migrationTask)
	if tools.PrintError(decodeErr) {
		return nil, false
	}
	return &migrationTask, true
}
