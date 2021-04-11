package model

/* get storage plan*/
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

/* get download plan*/
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

/* post storage plan*/
type PostStoragePlan struct {
	CloudID     string      `json:"CloudID" bson:"cloud_id"`
	UserID      string      `json:"UserId" bson:"user_id"`
	StoragePlan StoragePlan `json:"StoragePlan" bson:"storage_plan"`
}

type PostStoragePlanResponse struct {
	Code      uint64             `json:"Code" json:"code"`
	RequestID string             `json:"RequestID" bson:"request_id"`
	Data      []AccessCredential `json:"Data" bson:"data"`
	Msg       string             `json:"Msg" bson:"msg"`
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
