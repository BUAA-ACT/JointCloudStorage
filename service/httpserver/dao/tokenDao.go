package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func (d *Dao) InsertVerifyCode(email string, code string) bool {
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
	_, err := col.UpdateOne(context.TODO(), filter, update, opts)
	return !tools.PrintError(err)
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

func (d *Dao) InsertAccessToken(token string, userID string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	timeNow := time.Now()
	accessToken := &model.AccessTokenCode{
		AccessToken:           token,
		UserID:                userID,
		AccessTokenCreateTime: timeNow,
		AccessTokenModifyTime: timeNow,
	}
	_, err := col.InsertOne(context.TODO(), accessToken)
	return !tools.PrintError(err)
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
	// find it and check its time
	findErr := col.FindOne(context.TODO(), filter).Decode(&originToken)
	if findErr != nil {
		log.Println("无效token")
		return "", false
	}
	if timeNow.Sub(originToken.AccessTokenModifyTime).Hours() >= 4.0 || timeNow.Sub(originToken.AccessTokenCreateTime).Hours() >= 14*24.0 {
		// out of date and need to delete
		_, err := col.DeleteMany(context.TODO(), filter)
		tools.PrintError(err)
		return "", false
	} else {
		// in time and need to update
		_, err := col.UpdateMany(context.TODO(), filter, update)
		tools.PrintError(err)
		return originToken.UserID, true
	}
}

func (d *Dao) DeleteAccessToken(accessToken string) (*mongo.DeleteResult, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"access_token": accessToken,
	}
	result, err := col.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(err)
}
