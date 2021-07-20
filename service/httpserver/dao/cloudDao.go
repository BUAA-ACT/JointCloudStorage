package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	cloud request
*/
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

func (d *Dao) GetCloud(cloudID string) (*model.Cloud, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"cloud_id": cloudID,
	}
	var cloud model.Cloud
	findCloudError := col.FindOne(context.TODO(), filter).Decode(cloud)
	if tools.PrintError(findCloudError) {
		return nil, false
	}
	return &cloud, true
}

func (d *Dao) CheckSameCloudID(cloudID string) bool {
	filter := bson.M{
		"cloud_id": cloudID,
	}
	col := d.client.Database(d.database).Collection(d.collection)
	result := col.FindOne(context.TODO(), filter)
	return result.Err() == nil
}

func (d *Dao) GetAllVoteCloud() (*[]model.CloudController, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{}
	var cloudControllers = make([]model.CloudController, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var cloudController model.CloudController
		err := result.Decode(&cloudController)
		if tools.PrintError(err) {
			return nil, false
		}
		cloudControllers = append(cloudControllers, cloudController)
	}
	return &cloudControllers, true
}

func (d *Dao) GetAllAddedCloud() (*[]model.CloudController, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{}
	var addedClouds = make([]model.CloudController, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var addedCloud model.CloudController
		err := result.Decode(&addedCloud)
		if tools.PrintError(err) {
			return nil, false
		}
		addedClouds = append(addedClouds, addedCloud)
	}
	return &addedClouds, true
}
