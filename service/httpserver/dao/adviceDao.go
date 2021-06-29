package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Dao) GetNewAdvice(userID string) (*[]model.MigrationAdvice, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	var advices = make([]model.MigrationAdvice, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var advice model.MigrationAdvice
		err := result.Decode(&advice)
		if tools.PrintError(err) {
			return nil, false
		}
		if advice.CloudsNew == nil {
			advice.CloudsNew = make([]model.Cloud, 0)
		}
		if advice.CloudsOld == nil {
			advice.CloudsOld = make([]model.Cloud, 0)
		}
		advices = append(advices, advice)
	}
	return &advices, true
}

func (d *Dao) DeleteAdvice(userID string) (*mongo.DeleteResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	result, err := col.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(err)
}
