package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) GetNewAdvice(userId string) (*[]model.MigrationAdvice, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	var advices []model.MigrationAdvice
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
		advices = append(advices, advice)
	}
	return &advices, true
}

func (d *Dao) DeleteAdvice(userId string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	_, _ = col.DeleteMany(context.TODO(), filter)
}
