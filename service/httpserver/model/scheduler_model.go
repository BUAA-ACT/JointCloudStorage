package model

// GetStoragePlan data struct
type GetStoragePlan = Preference

type GetStoragePlanResponseData struct {
	StoragePriceFirst StoragePlan `json:"StoragePriceFirst" bson:"storage_price_first"`
	TrafficPriceFirst StoragePlan `json:"TrafficPriceFirst" bson:"traffic_price_first"`
}

type GetStoragePlanResponse struct {
	RequestID string                     `json:"RequestID" bson:"request_id"`
	Data      GetStoragePlanResponseData `json:"Data" bson:"data"`
	Code      uint64                     `json:"Code" bson:"code"`
	Msg       string                     `json:"Msg" bson:"msg"`
}

// GetDownloadPlan data struct
type GetDownloadPlan struct {
	UserID string `json:"UserID" bson:"user_id"`
	FileID string `json:"FileID" json:"file_id"`
}

type GetDownloadPlanResponseData = StoragePlan

type GetDownloadPlanResponse struct {
	Code      uint64                      `json:"Code" bson:"code"`
	RequestID string                      `json:"RequestID" bson:"request_id"`
	Data      GetDownloadPlanResponseData `json:"Data" bson:"data"`
	Msg       string                      `json:"Msg" bson:"msg"`
}

// PostStoragePlan data struct
type PostStoragePlan struct {
	CloudID     string      `json:"CloudID" bson:"cloud_id"`
	UserID      string      `json:"UserID" bson:"user_id"`
	Password    string      `json:"Password" bson:"password"`
	StoragePlan StoragePlan `json:"StoragePlan" bson:"storage_plan"`
}

type PostStoragePlanResponse struct {
	Code      uint64             `json:"Code" json:"code"`
	RequestID string             `json:"RequestID" bson:"request_id"`
	Data      []AccessCredential `json:"Data" bson:"data"`
	Msg       string             `json:"Msg" bson:"msg"`
}

// GetAllClouds data struct
type GetAllClouds struct {
}

type GetAllCloudsResponse struct {
	Code      uint64  `json:"Code" json:"code"`
	RequestID string  `json:"RequestID" bson:"request_id"`
	Data      []Cloud `json:"Data" bson:"data"`
	Msg       string  `json:"Msg" bson:"msg"`
}

// PostNewCloud data struct
type PostNewCloud = Cloud

type PostNewCloudResponse struct {
	Code      uint64 `json:"Code" json:"code"`
	RequestID string `json:"RequestID" bson:"request_id"`
	Msg       string `json:"Msg" bson:"msg"`
}

// PostUpdateCloud data struct
type PostUpdateCloud = Cloud

type PostUpdateCloudResponse struct {
	Code      uint64 `json:"Code" json:"code"`
	RequestID string `json:"RequestID" bson:"request_id"`
	Msg       string `json:"Msg" bson:"msg"`
}

// GetVoteRequests data struct
type GetVoteRequests struct {
}

type GetVoteRequestsResponse struct {
	Code      uint64  `json:"Code" json:"code"`
	RequestID string  `json:"RequestID" bson:"request_id"`
	Msg       string  `json:"Msg" bson:"msg"`
	Data      []Cloud `json:"Data" bson:"data"`
}

// PostCloudVote data struct
type PostCloudVote struct {
	CloudID    string `json:"CloudID" bson:"cloud_id"`
	VoteResult bool   `json:"VoteResult" bson:"vote_result"`
}

type PostCloudVoteResponse struct {
	Code      uint64 `json:"Code" json:"code"`
	RequestID string `json:"RequestID" bson:"request_id"`
	Msg       string `json:"Msg" bson:"msg"`
}

// PostKeyToScheduler data struct
type PostKeyToScheduler = AccessKey

type PostKeyToSchedulerResponse struct {
	Code      uint64 `json:"Code" json:"code"`
	RequestID string `json:"RequestID" bson:"request_id"`
	Msg       string `json:"Msg" bson:"msg"`
}

// DeleteKeyToScheduler data struct
type DeleteKeyToScheduler = AccessKey

type DeleteKeyToSchedulerResponse struct {
	Code      uint64 `json:"Code" json:"code"`
	RequestID string `json:"RequestID" bson:"request_id"`
	Msg       string `json:"Msg" bson:"msg"`
}

func (thisPlan *StoragePlan) isEqual(otherPlan *StoragePlan) bool {
	if thisPlan.StorageMode != otherPlan.StorageMode {
		return false
	}
	if thisPlan.N != otherPlan.N {
		return false
	}
	if thisPlan.K != otherPlan.K {
		return false
	}
	if thisPlan.StoragePrice != otherPlan.StoragePrice {
		return false
	}
	if thisPlan.TrafficPrice != otherPlan.TrafficPrice {
		return false
	}
	if thisPlan.Availability != otherPlan.Availability {
		return false
	}
	if len(thisPlan.Clouds) != len(otherPlan.Clouds) {
		return false
	}
	for index, cloud := range thisPlan.Clouds {
		if cloud != otherPlan.Clouds[index] {
			return false
		}
	}
	return true
}
