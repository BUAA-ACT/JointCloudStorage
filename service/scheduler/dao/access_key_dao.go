package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shaoliyin.me/jcspan/entity"
	"shaoliyin.me/jcspan/tools"
)

/*******************************************************
* 用于accesskey和secretKey的存储与删除
 ********************************************************/

//更新,若不存在则将插入
func KeyUpsert(col *mongo.Collection, ak entity.AccessKey) error {
	//col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	operation := bson.M{
		"$set": ak,
	}
	option := options.UpdateOptions{
		Upsert: tools.Bool2Pointer(true),
	}
	_, err := col.UpdateOne(context.TODO(), filter, operation, &option)
	if err != nil {
		return err
	}

	return nil
}

//删除key
func DeleteKey(col *mongo.Collection, ak entity.AccessKey) error {
	//col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	_, err := col.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}
	return nil
}
