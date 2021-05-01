package dao

import (
	"cloud-storage-httpserver/args"
	. "cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
	"cloud-storage-httpserver/service/tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) CreateNewUser(user User) {
	col := d.client.Database(d.database).Collection(d.collection)
	_, _ = col.InsertOne(context.TODO(), user)
}

func (d *Dao) GetUserInfo(userId string) (*User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var user User
	filter := bson.M{
		"user_id": userId,
	}
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	if user.AccessCredentials == nil {
		user.AccessCredentials = make([]AccessCredential, 0)
	}
	if user.DataStats.DownloadTraffic == nil {
		user.DataStats.DownloadTraffic = make(map[string]uint64)
		user.DataStats.UploadTraffic = make(map[string]uint64)
	}
	if user.Preference.Latency == nil {
		user.Preference.Latency = make(map[string]uint64)
	}
	if user.StoragePlan.Clouds == nil {
		user.StoragePlan.Clouds = make([]Cloud, 0)
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

func (d *Dao) SetUserStatusWithEmail(email string, status string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"email": email,
	}
	update := bson.D{{"$set",
		bson.D{
			{"user_status", status},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserStatusWithId(userId string, status string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"user_status", status},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) LoginWithEmail(email string, password string) (*User, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"email":    email,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user User
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	return &user, err == nil
}

func (d *Dao) LoginWithId(userId string, password string) bool {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id":  userId,
		"password": code.AesEncrypt(password, *args.EncryptKey),
	}
	var user User
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	return err == nil
}

func (d *Dao) SetUserPassword(userId string, password string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"password", code.AesEncrypt(password, *args.EncryptKey)},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserEmail(userId string, newEmail string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"email", newEmail},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserNickname(userId string, newNickname string) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"nickname", newNickname},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserPreference(userId string, preference *Preference) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"preference", *preference},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserStoragePlan(userId string, plan *StoragePlan) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"storage_plan", *plan},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}

func (d *Dao) SetUserAccessCredential(userId string, credentials *[]AccessCredential) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"user_id": userId,
	}
	update := bson.D{{"$set",
		bson.D{
			{"access_credentials", *credentials},
		},
	}}
	_, _ = col.UpdateMany(context.TODO(), filter, update)
}
