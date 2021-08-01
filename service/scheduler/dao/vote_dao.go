package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/entity"
)

/*
 * 下面函数用于操作投票类型voteCloud
 */

// InsertVoteCloud insert cloud into vote cloud collection
func (d *Dao) InsertVoteCloud(col *mongo.Collection, cloud entity.VoteCloud) error {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.InsertOne(context.TODO(), cloud)
	if err != nil {
		return err
	}
	return nil
}

// CloudsCount get the number of clouds whose id is cid
func (d *Dao) CloudsCount(col *mongo.Collection, cid string) (int64, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	count, err := col.CountDocuments(context.TODO(), bson.M{"cloud_id": cid})
	if err != nil {
		return count, err
	} else {
		return count, nil
	}
}

// DeleteVoteCloud delete the cloud
func (d *Dao) DeleteVoteCloud(col *mongo.Collection, id string) error {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	_, err := col.DeleteOne(context.TODO(), bson.M{"cloud_id": id})
	if err != nil {
		return err
	} else {
		return nil
	}
}

// AddVoteNum add vote number
func (d *Dao) AddVoteNum(col *mongo.Collection, vote int, id string) (int, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	res, err := col.UpdateOne(
		context.TODO(),
		bson.M{"cloud_id": id},
		bson.M{
			"$inc": bson.M{"vote_num": vote},
		})
	if err != nil && res != nil {
		return int(res.ModifiedCount), err
	} else if res != nil {
		return int(res.ModifiedCount), nil
	} else {
		return -1, err
	}
}

// GetVoteCloud Get struct voteCloud by id
func (d *Dao) GetVoteCloud(col *mongo.Collection, id string) (entity.VoteCloud, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result entity.VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

// GetAllVoteCloud Get all voteCloud in collection voteCloud
func (d *Dao) GetAllVoteCloud(col *mongo.Collection) ([]entity.VoteCloud, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result []entity.VoteCloud
	cur, err := col.Find(context.TODO(), bson.M{})
	defer cur.Close(context.TODO())
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var cloud entity.VoteCloud
		if err = cur.Decode(&cloud); err != nil {
			return result, err
		}
		result = append(result, cloud)
	}
	return result, nil
}

// GetVoteNumber : Get the vote number of the cloud with id
func (d *Dao) GetVoteNumber(col *mongo.Collection, id string) (int, error) {
	//col := d.client.Database(d.database).Collection(d.cloudCollection)
	var result entity.VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return -1, err
	} else {
		return result.VoteNum, nil
	}
}
