package model

type MigrationAdvice struct {
	UserID         string      `json:"UserID" bson:"user_id"`
	StoragePlanOld StoragePlan `json:"StoragePlanOld" bson:"storage_plan_old"`
	StoragePlanNew StoragePlan `json:"StoragePlanNew" bson:"storage_plan_new"`
	CloudsOld      []Cloud     `json:"CloudsOld" bson:"clouds_old"`
	CloudsNew      []Cloud     `json:"CloudsNew" bson:"clouds_new"`
	Cost           float64     `json:"Cost" json:"cost"`
	Status         string      `json:"Status" bson:"status"`
}
