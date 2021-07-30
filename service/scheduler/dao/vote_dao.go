package dao

import "go.mongodb.org/mongo-driver/bson"

/*
 * 下面函数用于操作投票类型voteCloud
 */

func (d *Dao) InsertVoteCloud(cloud VoteCloud) error {
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

//get the number of clouds whose id is cid
func (d *Dao) CloudsCount(cid string) (int64, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	count, err := col.CountDocuments(context.TODO(), bson.M{"cloud_id": cid})
	if err != nil {
		return count, err
	} else {
		return count, nil
	}
}

//delete the cloud
func (d *Dao) DeleteVoteCloud(id string) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	_, err := col.DeleteOne(context.TODO(), bson.M{"cloud_id": id})
	if err != nil {
		return err
	} else {
		return nil
	}
}

//add vote number
func (d *Dao) AddVoteNum(vote int, id string) (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	res, err := col.UpdateOne(
		context.TODO(),
		bson.M{"cloud_id": id},
		bson.M{
			"$inc": bson.M{"vote_num": vote},
		})
	if err != nil {
		return int(res.ModifiedCount), err
	} else {
		return int(res.ModifiedCount), nil
	}
}

//Get struct voteCloud by id
func (d *Dao) GetVoteCloud(id string) (VoteCloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

//Get all voteCloud in collection voteCloud
func (d *Dao) GetAllVoteCloud() ([]VoteCloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result []VoteCloud
	cur, err := col.Find(context.TODO(), bson.M{})
	defer cur.Close(context.TODO())
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var cloud VoteCloud
		if err = cur.Decode(&cloud); err != nil {
			return result, err
		}
		result = append(result, cloud)
	}
	return result, nil
}

//Get the vote number of the cloud with id
func (d *Dao) GetVoteNumber(id string) (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	var result VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return -1, err
	} else {
		return result.VoteNum, nil
	}
}
