package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (d *Dao) GetAllKeys(userID string) (*[]model.AccessKey, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	var keys = make([]model.AccessKey, 0)
	result, findErr := col.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var key model.AccessKey
		decodeErr := result.Decode(&key)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		keys = append(keys, key)
	}
	return &keys, true
}

func (d *Dao) GetWithAccessKey(userID string, accessKey string) (*model.AccessKey, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userID,
		"access_key": accessKey,
	}
	var key model.AccessKey
	findErr := col.FindOne(context.TODO(), filter).Decode(&key)
	if tools.PrintError(findErr) {
		return nil, false
	}
	return &key, true
}

func (d *Dao) InsertKey(userID string, accessKey string, secretKey string, comment string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	key := &model.AccessKey{
		UserID:     userID,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Comment:    comment,
		CreateTime: timeNow,
		Available:  true,
	}
	_, insertErr := col.InsertOne(context.TODO(), key)
	return !tools.PrintError(insertErr)
}

func (d *Dao) DeleteKey(userID string, accessKey string) (*mongo.DeleteResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userID,
		"access_key": accessKey,
	}
	result, deleteErr := col.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(deleteErr)
}

func (d *Dao) ChangeKeyStatus(userID string, accessKey string, status bool) (*mongo.UpdateResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userID,
		"access_key": accessKey,
	}
	update := bson.D{{"$set",
		bson.D{
			{"available", status},
		},
	}}
	result, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(changeErr)
}

func (d *Dao) RemakeKey(userID string, accessKey string, secretKey string) (*mongo.UpdateResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userID,
		"access_key": accessKey,
	}
	update := bson.D{{"$set",
		bson.D{
			{"secret_key", secretKey},
		},
	}}
	result, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(changeErr)
}

func (d *Dao) ChangeKeyComment(userID string, accessKey string, newComment string) (*mongo.UpdateResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userID,
		"access_key": accessKey,
	}
	update := bson.D{{"$set",
		bson.D{
			{"comment", newComment},
		},
	}}
	result, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(changeErr)
}
