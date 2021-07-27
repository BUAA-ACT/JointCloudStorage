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
	filter := bson.M{}
	var clouds = make([]model.Cloud, 0)
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var cloud model.Cloud
		decodeErr := result.Decode(&cloud)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		clouds = append(clouds, cloud)
	}
	return &clouds, true
}

func (d *Dao) GetCloud(cloudID string) (*model.Cloud, bool) {
	filter := bson.M{
		"cloud_id": cloudID,
	}
	var cloud model.Cloud
	findCloudError := d.collectionConnection.FindOne(context.TODO(), filter).Decode(&cloud)
	if tools.PrintError(findCloudError) {
		return nil, false
	}
	return &cloud, true
}

func (d *Dao) CheckSameCloudID(cloudID string) bool {
	filter := bson.M{
		"cloud_id": cloudID,
	}
	result := d.collectionConnection.FindOne(context.TODO(), filter)
	return result.Err() == nil
}

func (d *Dao) GetAllVoteCloud() (*[]model.CloudController, bool) {
	filter := bson.M{}
	var cloudControllers = make([]model.CloudController, 0)
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var cloudController model.CloudController
		decodeErr := result.Decode(&cloudController)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		cloudControllers = append(cloudControllers, cloudController)
	}
	return &cloudControllers, true
}

func (d *Dao) GetAllAddedCloud() (*[]model.CloudController, bool) {
	filter := bson.M{}
	var addedClouds = make([]model.CloudController, 0)
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var addedCloud model.CloudController
		decodeErr := result.Decode(&addedCloud)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		addedClouds = append(addedClouds, addedCloud)
	}
	return &addedClouds, true
}
