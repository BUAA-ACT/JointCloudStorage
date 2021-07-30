package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"shaoliyin.me/jcspan/entity"
)

// UpdateCloud insert new cloud info to database.
func (d *Dao) UpdateCloud(cloud entity.Cloud) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
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
func (d *Dao) GetAllClouds() ([]entity.Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

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

func (d *Dao) GetOtherClouds(cid string) ([]entity.Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

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

func (d *Dao) GetCloud(cid string) (entity.Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var cloud entity.Cloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": cid}).Decode(&cloud)

	// 隐藏访问凭证
	cloud.AccessKey = ""
	cloud.SecretKey = ""
	return cloud, err
}

func (d *Dao) GetCloudNum() (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	num, err := col.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	} else {
		return int(num), nil
	}
}

func (d *Dao) InsertCloud(cloud entity.Cloud) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.InsertOne(
		context.TODO(),
		cloud,
	)
	if err != nil {
		return err
	}
	return nil
}
