package entity

type BaseResponse struct {
	RequestID string
	Code      int
	Msg       string
}

type GetStoragePlanParam Preference

type GetStoragePlanData struct {
	StoragePriceFirst StoragePlan
	TrafficPriceFirst StoragePlan
}

type GetDownloadPlanParam struct {
	UserID string
	FileID string
}

type GetDownloadPlanData struct {
	StorageMode string
	Clouds      []Cloud
	Index       []int
}

type GetStatusParam struct {
	CloudID string
}

type GetStatusData struct {
	Cloud
}

type PostStoragePlanParam struct {
	CloudID     string
	UserID      string
	Password    string
	StoragePlan StoragePlan
}

type PostStoragePlanData struct {
	AccessCredential
}

type PostMetadataParam struct {
	CloudID string
	UserID  string
	Type    string
	Clouds  []Cloud
	Files   []File
}

type PostMetadataData struct {
}

type GetStatusResponse struct {
	BaseResponse
	Data Cloud
}

type PostStoragePlanResponse struct {
	BaseResponse
	Data []AccessCredential
}

type PostMetadataResponse struct {
	BaseResponse
}
