package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*******************************************************
* 用于accesskey和secretKey的存储与删除
 ********************************************************/

//更新,若不存在则将插入
func (dao *Dao) KeyUpsert(ak AccessKey) error {
	col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	operation := bson.M{
		"$set": ak,
	}
	option := options.UpdateOptions{
		Upsert: bool2pointer(true),
	}
	_, err := col.UpdateOne(context.TODO(), filter, operation, &option)
	if err != nil {
		return err
	}

	return nil
}

//删除key
func (dao *Dao) DeleteKey(ak AccessKey) error {
	col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	_, err := col.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}
	return nil
}
