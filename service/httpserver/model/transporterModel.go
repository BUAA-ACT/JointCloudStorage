package model

type TaskResponseData struct {
	Type   string `json:"Type" bson:"type"`
	Result string `json:"Result" bson:"result"`
}

type TaskResponse struct {
	Code uint64           `json:"Code" bson:"code"`
	Msg  string           `json:"Msg" bson:"msg"`
	Data TaskResponseData `json:"Data" bson:"data"`
}
