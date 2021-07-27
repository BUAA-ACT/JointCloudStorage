package dao

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var count = 0
var sum1 int64 = 0
var sum2 int64 = 0

func (d *Dao) InsertVerifyCode(email string, code string) bool {
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
	_, insertErr := d.collectionConnection.UpdateOne(context.TODO(), filter, update, opts)
	return !tools.PrintError(insertErr)
}

func (d *Dao) VerifyEmail(email string, verifyCode string) bool {
	filter := bson.M{
		"email":       email,
		"verify_code": verifyCode,
	}
	var verifyCodeResult model.VerifyCode
	parseJsonError := d.collectionConnection.FindOneAndDelete(context.TODO(), filter).Decode(&verifyCodeResult)
	if parseJsonError != nil {
		return false
	}
	timeNow := time.Now()
	if timeNow.Sub(verifyCodeResult.VerifyCodeCreateTime) >= args.VerifyCodeTimeOut {
		return false
	}
	return true
}

func (d *Dao) InsertAccessToken(token string, userID string) bool {
	timeNow := time.Now()
	accessToken := &model.AccessTokenCode{
		AccessToken:           token,
		UserID:                userID,
		AccessTokenCreateTime: timeNow,
		AccessTokenModifyTime: timeNow,
	}
	_, insertErr := d.collectionConnection.InsertOne(context.TODO(), accessToken)
	return !tools.PrintError(insertErr)
}

func (d *Dao) CheckValid(token string) (string, bool) {
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
	findErr := d.collectionConnection.FindOne(context.TODO(), filter).Decode(&originToken)
	if findErr != nil {
		log.Println("无效token")
		return "", false
	}
	if timeNow.Sub(originToken.AccessTokenModifyTime) >= args.AccessTokenTimeOut || timeNow.Sub(originToken.AccessTokenCreateTime) >= args.AccessTokenDateOut {
		// out of date and need to delete
		_, deleteErr := d.collectionConnection.DeleteMany(context.TODO(), filter)
		tools.PrintError(deleteErr)
		return "", false
	} else {
		// in time and need to update
		_, changeErr := d.collectionConnection.UpdateMany(context.TODO(), filter, update)
		tools.PrintError(changeErr)
		return originToken.UserID, true
	}
}

func (d *Dao) DeleteAccessToken(accessToken string) (*mongo.DeleteResult, bool) {
	filter := bson.M{
		"access_token": accessToken,
	}
	result, deleteErr := d.collectionConnection.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(deleteErr)
}

func (d *Dao) CleanAccessToken() (*mongo.DeleteResult, bool) {
	timeNow := time.Now()
	accessTokenModifyDeadLine := timeNow.Add(-args.AccessTokenTimeOut)
	accessTokenCreateDeadLine := timeNow.Add(-args.AccessTokenDateOut)
	filter := bson.M{
		"$or": []bson.M{
			{
				"$or": []bson.M{
					{
						"access_token_create_time": bson.M{
							"$gt": timeNow,
						},
					},
					{
						"access_token_create_time": bson.M{
							"lt": accessTokenCreateDeadLine,
						},
					},
				},
			},
			{
				"$or": []bson.M{
					{
						"access_token_modify_time": bson.M{
							"$gt": timeNow,
						},
					},
					{
						"access_token_modify_time": bson.M{
							"$lt": accessTokenModifyDeadLine,
						},
					},
				},
			},
		},
	}
	result, deleteErr := d.collectionConnection.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(deleteErr)
}

func (d *Dao) CleanVerifyCode() (*mongo.DeleteResult, bool) {
	timeNow := time.Now()
	verifyCodeDeadLine := timeNow.Add(-args.VerifyCodeTimeOut)
	filter := bson.M{
		"$or": []bson.M{
			{
				"verify_code_create_time": bson.M{
					"$gt": timeNow,
				},
			},
			{
				"verify_code_create_time": bson.M{
					"$lt": verifyCodeDeadLine,
				},
			},
		},
	}
	result, deleteErr := d.collectionConnection.DeleteMany(context.TODO(), filter)
	return result, !tools.PrintError(deleteErr)
}
