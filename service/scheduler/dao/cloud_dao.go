package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/entity"
)

// UpdateCloud insert new cloud info to database.
func UpdateCloud(col *mongo.Collection, cloud entity.Cloud) error {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"cloud_id": cloud.CloudID,
		},
		bson.M{
			"$set": bson.M{
				"storage_price": cloud.StoragePrice,
				"traffic_price": cloud.TrafficPrice,
				"availability":  cloud.Availability,
				"status":        cloud.Status,
				"location":      cloud.Location,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAllClouds return the info of given bucket.
func GetAllClouds(col *mongo.Collection) ([]entity.Cloud, error) {
	//col := d.Client.Database(d.database).Collection(d.cloudCollection)
	var clouds []entity.Cloud
	cur, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem entity.Cloud
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		// 隐藏访问凭证
		//elem.AccessKey = ""
		//elem.SecretKey = ""
		clouds = append(clouds, elem)
	}

	return clouds, nil
}

func GetOtherClouds(col *mongo.Collection, cid string) ([]entity.Cloud, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	var clouds []entity.Cloud
	cur, err := col.Find(context.TODO(), bson.M{"cloud_id": bson.M{"$ne": cid}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem entity.Cloud
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		// 隐藏访问凭证
		elem.AccessKey = ""
		elem.SecretKey = ""
		clouds = append(clouds, elem)
	}

	return clouds, nil
}

func GetCloud(col *mongo.Collection, cid string) (entity.Cloud, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	var cloud entity.Cloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": cid}).Decode(&cloud)

	// 隐藏访问凭证
	cloud.AccessKey = ""
	cloud.SecretKey = ""
	return cloud, err
}

func GetCloudNum(col *mongo.Collection) (int, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	num, err := col.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	} else {
		return int(num), nil
	}
}

func InsertCloud(col *mongo.Collection, cloud entity.Cloud) error {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.InsertOne(
		context.TODO(),
		cloud,
	)
	if err != nil {
		return err
	}
	return nil
}
