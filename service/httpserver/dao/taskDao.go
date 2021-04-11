package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) GetTask(taskId string, userId string, isSingle bool) (*[]model.Task, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var filter interface{}
	filterSingle := bson.M{
		"task_id": taskId,
	}
	filterAll := bson.M{
		"user_id": userId,
	}
	if isSingle {
		filter = filterSingle
	} else {
		filter = filterAll
	}
	var tasks []model.Task
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var task model.Task
		err := result.Decode(&task)
		if tools.PrintError(err) {
			return nil, false
		}
		tasks = append(tasks, task)
	}
	return &tasks, true
}
