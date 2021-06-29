package model

import "time"

type AccessKey struct {
	UserID     string    `json:"UserID" bson:"user_id"`
	AccessKey  string    `json:"AccessKey" bson:"access_key"`
	SecretKey  string    `json:"SecretKey" bson:"secret_key"`
	Comment    string    `json:"Comment" bson:"comment"`
	CreateTime time.Time `json:"CreateTime" bson:"create_time"`
	Available  bool      `json:"Available" bson:"available"`
}
