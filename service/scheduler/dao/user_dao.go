package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Dao) GetUser(uid string) (User, error) {
	col := d.client.Database(d.database).Collection(d.userCollection)

	var user User
	err := col.FindOne(context.TODO(), bson.M{"user_id": uid}).Decode(&user)
	return user, err
}

func (d *Dao) GetAllUser() ([]User, error) {
	col := d.client.Database(d.database).Collection(d.userCollection)

	var users []User
	cur, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		users = append(users, elem)
	}

	return users, nil
}

func (d *Dao) InsertUser(user User) error {
	col := d.client.Database(d.database).Collection(d.userCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"user_id": user.UserId,
		},
		bson.M{
			"$set": user,
		},
		&options.UpdateOptions{
			Upsert: bool2pointer(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) ChangeVolume(uid string, op string, files []File) error {
	var sum int64
	for _, v := range files {
		sum += v.Size
	}
	if op == "Delete" {
		sum = -sum
	}

	col := d.client.Database(d.database).Collection(d.userCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"user_id": uid,
		},
		bson.M{
			"$inc": bson.M{"data_stats.volume": sum},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) DeleteUser(uid string) error {
	// 删除该用户名下所有文件
	col := d.client.Database(d.database).Collection(d.fileCollection)
	_, err := col.DeleteMany(
		context.TODO(),
		bson.M{
			"owner": uid,
		},
	)
	if err != nil {
		return err
	}

	// 删除用户
	col = d.client.Database(d.database).Collection(d.userCollection)
	_, err = col.DeleteOne(
		context.TODO(),
		bson.M{
			"user_id": uid,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
