package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"time"
)

type User struct {
	UserId            string           `bson:"user_id"`
	Email             string           `bson:"email"`
	Password          string           `bson:"password"`
	Nickname          string           `bson:"nickname"`
	Role              string           `bson:"role"`
	Avatar            string           `bson:"avatar"`
	LastModified      time.Time        `bson:"last_modified"`
	Preference        Preference       `bson:"preference"`
	StoragePlan       StoragePlan      `bson:"storage_plan"`
	DataStats         DataStats        `bson:"data_stats"`
	AccessCredentials AccessCredential `bson:"access_credentials"`
	Status            string           `bson:"status"`
}

type Preference struct {
	Vendor       int            `bson:"vendor"`
	StoragePrice float64        `bson:"storage_price"`
	TrafficPrice float64        `bson:"traffic_price"`
	Availability float64        `bson:"availability"`
	Latency      map[string]int `bson:"latency"`
}

type DataStats struct {
	Volume          int64            `bson:"volume"`
	UploadTraffic   map[string]int64 `bson:"upload_traffic"`
	DownloadTraffic map[string]int64 `bson:"download_traffic"`
}

type AccessCredential struct {
	CloudID  string `bson:"cloud_id"`
	UserID   string `bson:"user_id"`
	Password string `bson:"password"`
}

type UserDatabase interface {
	GetUserFromID(uid string) (*User, error)
	UpdateUserInfo(user *User) error
}

type InMemoryUserDatabase struct {
	UserInfoMap map[string]User
}

func NewInMemoryUserDatabase() *InMemoryUserDatabase {
	return &InMemoryUserDatabase{UserInfoMap: map[string]User{
		"tester": {
			UserId:       "tester",
			Email:        "",
			Password:     "",
			Nickname:     "",
			Role:         "",
			Avatar:       "",
			LastModified: time.Time{},
			Preference:   Preference{},
			StoragePlan:  StoragePlan{},
			DataStats: DataStats{
				Volume: 0,
				UploadTraffic: map[string]int64{
					util.Config.LocalCloudID: 0,
				},
				DownloadTraffic: map[string]int64{
					util.Config.LocalCloudID: 0,
				},
			},
			AccessCredentials: AccessCredential{},
			Status:            "",
		},
	}}
}

func (db *InMemoryUserDatabase) GetUserFromID(uid string) (*User, error) {
	user := db.UserInfoMap[uid]
	return &user, nil
}
func (db *InMemoryUserDatabase) UpdateUserInfo(user *User) error {
	db.UserInfoMap[user.UserId] = *user
	return nil
}
