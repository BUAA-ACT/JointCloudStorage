package dao

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Dao) GetNewAdvice(userID string) (*[]model.MigrationAdvice, bool) {
	filter := bson.M{
		"user_id": userID,
	}
	var advices = make([]model.MigrationAdvice, 0)
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var advice model.MigrationAdvice
		decodeErr := result.Decode(&advice)
		if tools.PrintError(decodeErr) {
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
	filter := bson.M{
		"user_id": userID,
	}
	result, deleteErr := d.collectionConnection.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(deleteErr)
}

func (d *Dao) SetAdviceStatus(userID string, status string) (*mongo.UpdateResult, bool) {
	filter := bson.M{
		"user_id": userID,
		"status":  args.AdviceStatusPending,
	}
	update := bson.D{{"$set",
		bson.D{
			{"status", status},
		},
	}}
	result, changeErr := d.collectionConnection.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(changeErr)
}
