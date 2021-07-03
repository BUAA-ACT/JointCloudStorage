package model

import (
	"cloud-storage-httpserver/args"
	"time"
)

type Cloud struct {
	CloudID      string  `json:"CloudID,omitempty" bson:"cloud_id"`
	Endpoint     string  `json:"Endpoint,omitempty" bson:"endpoint"`
	AccessKey    string  `json:"AccessKey,omitempty" bson:"access_key"`
	SecretKey    string  `json:"SecretKey,omitempty" bson:"secret_key"`
	StoragePrice float64 `json:"StoragePrice" bson:"storage_price"`
	TrafficPrice float64 `json:"TrafficPrice" bson:"traffic_price"`
	Availability float64 `json:"Availability" bson:"availability"`
	Status       string  `json:"Status" bson:"status"`     // "UP" | "DOWN"
	Location     string  `json:"Location" bson:"location"` // "116.381252,39.906569"
	Address      string  `json:"Address" bson:"address"`
	CloudName    string  `json:"CloudName" bson:"cloud_name"`
	ProviderName string  `json:"ProviderName" bson:"provider_name"`
	Bucket       string  `json:"Bucket" bson:"bucket"`
}

type Preference struct {
	Vendor       uint64            `json:"Vendor" binding:"gt=0" bson:"vendor"`
	StoragePrice float64           `json:"StoragePrice" binding:"gt=0" bson:"storage_price"`
	TrafficPrice float64           `json:"TrafficPrice" binding:"gt=0" bson:"traffic_price"`
	Availability float64           `json:"Availability" binding:"gte=0,lte=1" bson:"availability"`
	Latency      map[string]uint64 `json:"Latency" bson:"latency"`
}

type StoragePlan struct {
	StorageMode  string  `json:"StorageMode" bson:"storage_mode"`
	N            int     `json:"N" bson:"n"`
	K            int     `json:"K" bson:"k"`
	Clouds       []Cloud `json:"Clouds" bson:"clouds"`
	StoragePrice float64 `json:"StoragePrice" bson:"storage_price"`
	TrafficPrice float64 `json:"TrafficPrice" bson:"traffic_price"`
	Availability float64 `json:"Availability" bson:"availability"`
}

type DataStats struct {
	Volume          uint64            `json:"Volume" bson:"volume"`
	UploadTraffic   map[string]uint64 `json:"UploadTraffic" bson:"upload_traffic"`
	DownloadTraffic map[string]uint64 `json:"DownloadTraffic" bson:"download_traffic"`
}

type AccessCredential struct {
	CloudID  string `json:"CloudID" bson:"cloud_id"`
	UserID   string `json:"UserID" bson:"user_id"`
	Password string `json:"Password" bson:"password"`
}

type User struct {
	UserID            string             `json:"UserID" bson:"user_id"`
	Email             string             `json:"Email" bson:"email"`
	Password          string             `json:"-" bson:"password"`
	Nickname          string             `json:"Nickname" bson:"nickname"`
	Role              string             `json:"Role" bson:"role"`
	Avatar            string             `json:"Avatar" bson:"avatar"`
	CreateTime        time.Time          `json:"CreateTime" bson:"create_time"`
	LastModified      time.Time          `json:"LastModified" bson:"last_modified"`
	Preference        Preference         `json:"Preference" bson:"preference"`
	StoragePlan       StoragePlan        `json:"StoragePlan" bson:"storage_plan"`
	DataStats         DataStats          `json:"DataStats" bson:"data_stats"`
	AccessCredentials []AccessCredential `json:"AccessCredentials" bson:"access_credentials"`
	Status            string             `json:"Status" bson:"status"`
	//TODO
}

func (user *User) UserHaveStoragePlan() bool {
	return user.StoragePlan.N > 0
}

func (user *User) UserHavePreference() bool {
	return user.Preference.Vendor > 0
}

func (user *User) IsNormalStatus() bool {
	return user.Role == args.UserNormalStatus
}

func (user *User) IsForbiddenStatus() bool {
	return user.Role == args.UserForbiddenStatus
}

func (user *User) IsVerifying() bool {
	return user.Role == args.UserVerifyStatus
}
