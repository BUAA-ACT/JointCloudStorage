package dao

import (
	"cloud-storage-httpserver/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (d *Dao) InsertVerifyCode(email string, code string) {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	upsert := true
	opts := &options.UpdateOptions{Upsert: &upsert}
	filter := bson.M{
		"email": email,
	}
	update := bson.D{{"$set",
		bson.D{
			{"verify_code", code},
			{"verify_code_create_time", timeNow},
		},
	}}
	_, _ = col.UpdateOne(context.TODO(), filter, update, opts)
}

func (d *Dao) VerifyEmail(email string, verifyCode string) bool {
	col := d.client.Database(d.database).Collection(d.collection)

	filter := bson.M{
		"email":       email,
		"verify_code": verifyCode,
	}
	result := col.FindOneAndDelete(context.TODO(), filter)
	if result.Err() != nil {
		return false
	}
	return true
}

func (d *Dao) InsertAccessToken(token string, userId string) {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	accessToken := &model.AccessTokenCode{
		AccessToken:           token,
		UserId:                userId,
		AccessTokenCreateTime: timeNow,
		AccessTokenModifyTime: timeNow,
	}
	_, _ = col.InsertOne(context.TODO(), accessToken)
}

func (d *Dao) CheckValid(token string) (string, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	filter := bson.M{
		"access_token": token,
	}
	update := bson.D{{"$set",
		bson.D{
			{"access_token_modify_time", timeNow},
		},
	}}
	var originToken model.AccessTokenCode
	err := col.FindOneAndUpdate(context.TODO(), filter, update).Decode(&originToken)
	if err != nil {
		fmt.Print(err)
		return "", false
	}
	return originToken.UserId, true
}

func (d *Dao) DeleteAccessToken(accessToken string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"access_token": accessToken,
	}
	_, _ = col.DeleteMany(context.TODO(), filter)
}
