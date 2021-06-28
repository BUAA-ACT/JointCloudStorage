package dao

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) CreateNewUser(user model.User) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	_, err := col.InsertOne(context.TODO(), user)
	return !tools.PrintError(err)
}

func (d *Dao) GetUserInfo(userId string) (*model.User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var user model.User
	filter := bson.M{
		"user_id": userId,
	}
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	if user.AccessCredentials == nil {
		user.AccessCredentials = make([]model.AccessCredential, 0)
	}
	if user.DataStats.DownloadTraffic == nil {
		user.DataStats.DownloadTraffic = make(map[string]uint64)
		user.DataStats.UploadTraffic = make(map[string]uint64)
	}
	if user.Preference.Latency == nil {
		user.Preference.Latency = make(map[string]uint64)
	}
	if user.StoragePlan.Clouds == nil {
		user.StoragePlan.Clouds = make([]model.Cloud, 0)
	}
	if tools.PrintError(err) {
		return nil, false
	}
	return &user, true
}

func (d *Dao) CheckSameEmail(email string) bool {
	filter := bson.M{
		"email": email,
	}
	col := d.client.Database(d.database).Collection(d.collection)
	result := col.FindOne(context.TODO(), filter)
	return !tools.PrintError(result.Err())
}

func (d *Dao) SetUserStatusWithEmail(email string, status string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"email": email,
	}
	update := bson.D{{"$set",
		bson.D{
			{"user_status", status},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserStatusWithId(userId string, status string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"user_status", status},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) LoginWithEmail(email string, password string) (*model.User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"email":    email,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user model.User
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	return &user, err == nil
}

func (d *Dao) LoginWithId(userId string, password string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":  userId,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user model.User
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	return err == nil
}

func (d *Dao) SetUserPassword(userId string, password string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"password", code.AesEncrypt(password, *args.EncryptKey)},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserEmail(userId string, newEmail string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"email", newEmail},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserNickname(userId string, newNickname string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"nickname", newNickname},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserPreference(userId string, preference *model.Preference) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"preference", *preference},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserStoragePlan(userId string, plan *model.StoragePlan) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"storage_plan", *plan},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}

func (d *Dao) SetUserAccessCredential(userId string, credentials *[]model.AccessCredential) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"access_credentials", *credentials},
		},
	}}
	_, err := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(err)
}
