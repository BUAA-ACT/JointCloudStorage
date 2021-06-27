package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

func (d *Dao) DeleteKey(userId string, accessKey string) {

}

func (d *Dao) ChangeKeyStatus(userId string, accessKey string) {

}

func (d *Dao) RemakeKey(userId string, accessKey string) {

}
