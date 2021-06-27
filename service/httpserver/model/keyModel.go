package model

import "time"

type AccessKey struct {
	UserId     string    `bson:"user_id"`
	AccessKey  string    `bson:"access_key"`
	SecretKey  string    `bson:"secret_key"`
	CreateTime time.Time `bson:"create_time"`
	Available  bool      `bson:"available"`
}
