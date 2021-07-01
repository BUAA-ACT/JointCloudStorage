package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) GetAllClouds() (*[]model.Cloud, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{}
	var clouds = make([]model.Cloud, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var cloud model.Cloud
		err := result.Decode(&cloud)
		if tools.PrintError(err) {
			return nil, false
		}
		clouds = append(clouds, cloud)
	}
	return &clouds, true
}

func (d *Dao) CheckSameCloudID(cloudID string) bool {
	filter := bson.M{
		"cloud_id": cloudID,
	}
	col := d.client.Database(d.database).Collection(d.collection)
	result := col.FindOne(context.TODO(), filter)
	return result.Err() == nil
}
