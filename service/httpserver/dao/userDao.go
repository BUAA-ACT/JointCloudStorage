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
	_, insertErr := col.InsertOne(context.TODO(), user)
	return !tools.PrintError(insertErr)
}

func (d *Dao) GetUserInfo(userID string) (*model.User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var user model.User
	filter := bson.M{
		"user_id": userID,
	}
	findErr := col.FindOne(context.TODO(), filter).Decode(&user)
	if user.AccessCredentials == nil {
		user.AccessCredentials = make([]model.AccessCredential, 0)
	}
	if user.DataStats.DownloadTraffic == nil {
		user.DataStats.DownloadTraffic = make(map[string]uint64)
	}
	if user.DataStats.UploadTraffic == nil {
		user.DataStats.UploadTraffic = make(map[string]uint64)
	}
	if user.Preference.Latency == nil {
		user.Preference.Latency = make(map[string]uint64)
	}
	if user.StoragePlan.Clouds == nil {
		user.StoragePlan.Clouds = make([]model.Cloud, 0)
	}
	if tools.PrintError(findErr) {
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
	return result.Err() == nil
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
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserStatusWithId(userID string, status string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"user_status", status},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) LoginWithEmail(email string, password string) (*model.User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"email":    email,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user model.User
	changeErr := col.FindOne(context.TODO(), filter).Decode(&user)
	return &user, changeErr == nil
}

func (d *Dao) LoginWithId(userID string, password string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":  userID,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user model.User
	findErr := col.FindOne(context.TODO(), filter).Decode(&user)
	return findErr == nil
}

func (d *Dao) SetUserPassword(userID string, password string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"password", code.AesEncrypt(password, *args.EncryptKey)},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserEmail(userID string, newEmail string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"email", newEmail},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserNickname(userID string, newNickname string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"nickname", newNickname},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserPreference(userID string, preference *model.Preference) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"preference", *preference},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserStoragePlan(userID string, plan *model.StoragePlan) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"storage_plan", *plan},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}

func (d *Dao) SetUserAccessCredential(userID string, credentials *[]model.AccessCredential) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.D{{"$set",
		bson.D{
			{"access_credentials", *credentials},
		},
	}}
	_, changeErr := col.UpdateMany(context.TODO(), filter, update)
	return !tools.PrintError(changeErr)
}
