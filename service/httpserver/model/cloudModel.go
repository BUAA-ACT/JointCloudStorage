package model

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

type CloudController struct {
	CloudID string `bson:"cloud_id" json:"CloudID"`
	Cloud   Cloud  `bson:"cloud" json:"Cloud"`
	VoteNum int    `bson:"vote_num" json:"VoteNum"`
	Address string `bson:"address" json:"Address"`
}
