package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) GetAllClouds(userId string) (*[]model.Cloud, bool) {
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
