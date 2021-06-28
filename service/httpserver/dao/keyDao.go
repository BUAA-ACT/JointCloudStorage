package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (d *Dao) GetAllKeys(userId string) (*[]model.AccessKey, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{}
	var keys = make([]model.AccessKey, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var key model.AccessKey
		err := result.Decode(&key)
		if tools.PrintError(err) {
			return nil, false
		}
		keys = append(keys, key)
	}
	return &keys, true
}

func (d *Dao) InsertKey(userId string, accessKey string, secretKey string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	key := &model.AccessKey{
		UserId:     userId,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		CreateTime: timeNow,
		Available:  true,
	}
	id, err := col.InsertOne(context.TODO(), key)
	fmt.Println("insert a key with id:", id)
	return !tools.PrintError(err)
}

func (d *Dao) DeleteKey(userId string, accessKey string) (*mongo.DeleteResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userId,
		"access_key": accessKey,
	}
	result, err := col.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(err)
}

func (d *Dao) ChangeKeyStatus(userId string, accessKey string, status bool) (*mongo.UpdateResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userId,
		"access_key": accessKey,
	}
	update := bson.D{{"$set",
		bson.D{
			{"available", status},
		},
	}}
	result, err := col.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(err)
}

func (d *Dao) RemakeKey(userId string, accessKey string, secretKey string) (*mongo.UpdateResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":    userId,
		"access_key": accessKey,
	}
	update := bson.D{{"$set",
		bson.D{
			{"secret_key", secretKey},
		},
	}}
	result, err := col.UpdateMany(context.TODO(), filter, update)
	return result, !tools.PrintError(err)
}
