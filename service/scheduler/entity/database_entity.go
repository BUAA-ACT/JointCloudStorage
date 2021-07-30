package entity

import "time"

type Cloud struct {
	CloudID      string  `bson:"cloud_id"`
	Endpoint     string  `bson:"endpoint"`
	AccessKey    string  `bson:"access_key" `
	SecretKey    string  `bson:"secret_key" `
	StoragePrice float64 `bson:"storage_price"`
	TrafficPrice float64 `bson:"traffic_price"`
	Availability float64 `bson:"availability"`
	Status       string  `bson:"status"`
	Location     string  `bson:"location"`
	Address      string  `bson:"address"`
	CloudName    string  `bson:"cloud_name"`
	ProviderName string  `bson:"provider_name"`
	Bucket       string  `json:"Bucket" bson:"bucket"`
}

type User struct {
	UserId            string             `bson:"user_id"`
	Email             string             `bson:"email"`
	Password          string             `bson:"password"`
	Nickname          string             `bson:"nickname"`
	Role              string             `bson:"role"`
	Avatar            string             `bson:"avatar"`
	LastModified      time.Time          `bson:"last_modified"`
	Preference        Preference         `bson:"preference"`
	StoragePlan       StoragePlan        `bson:"storage_plan"`
	DataStats         DataStats          `bson:"data_stats"`
	AccessCredentials []AccessCredential `bson:"access_credentials"`
	Status            string             `bson:"status"`
}

type Preference struct {
	Vendor       int            `bson:"vendor"`
	StoragePrice float64        `bson:"storage_price"`
	TrafficPrice float64        `bson:"traffic_price"`
	Availability float64        `bson:"availability"`
	Latency      map[string]int `bson:"latency"`
}

type StoragePlan struct {
	N            int     `bson:"n"`
	K            int     `bson:"k"`
	StorageMode  string  `bson:"storage_mode"`
	Clouds       []Cloud `bson:"clouds"`
	StoragePrice float64 `bson:"storage_price"`
	TrafficPrice float64 `bson:"traffic_price"`
	Availability float64 `bson:"availability"`
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

type File struct {
	FileID            string    `bson:"file_id"`
	FileName          string    `bson:"file_name"`
	Owner             string    `bson:"owner"`
	Size              int64     `bson:"size"`
	LastModified      time.Time `bson:"last_modified"`
	SyncStatus        string    `bson:"sync_status"`
	LastReconstructed time.Time `bson:"last_reconstructed"`
	ReconstructStatus string    `bson:"reconstruct_status"`
	DownloadUrl       string    `bson:"download_url"`
}

type MigrationAdvice struct {
	UserId         string      `bson:"user_id"`
	StoragePlanOld StoragePlan `bson:"storage_plan_old"`
	StoragePlanNew StoragePlan `bson:"storage_plan_new"`
	CloudsOld      []Cloud     `bson:"clouds_old"`
	CloudsNew      []Cloud     `bson:"clouds_new"`
	Cost           float64     `bson:"cost"`
	Status         string      `bson:"status"`
}

type AccessKey struct {
	UserID     string    `json:"UserID,omitempty" bson:"user_id"`
	AccessKey  string    `json:"AccessKey,omitempty" bson:"access_key"`
	SecretKey  string    `json:"SecretKey,omitempty" bson:"secret_key"`
	Comment    string    `json:"Comment,omitempty" bson:"comment"`
	CreateTime time.Time `json:"CreateTime,omitempty" bson:"create_time"`
	Available  bool      `json:"Available,omitempty" bson:"available"`
}

type VoteCloud struct {
	Id      string `bson:"cloud_id" json:"id"`
	Cloud   Cloud  `bson:"cloud" json:"cloud"`
	VoteNum int    `bson:"vote_num" json:"vote_num"`
	Address string `bson:"address" json:"address"`
}
