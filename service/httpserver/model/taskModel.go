package model

import "time"

type Task struct {
	TaskID                 string      `json:"TaskID" bson:"task_id"`
	TaskType               string      `json:"TaskType" bson:"task_type"`
	TaskState              string      `json:"TaskState" bson:"task_state"`
	TaskStartTime          time.Time   `json:"TaskStartTime" bson:"task_start_time"`
	UserId                 string      `json:"UserID" bson:"user_id"`
	SourcePath             string      `json:"SourcePath,omitempty" bson:"source_path"`
	DestinationPath        string      `json:"DestinationPath,omitempty" bson:"destination_path"`
	SourceStoragePlan      StoragePlan `json:"SourceStoragePlan,omitempty" bson:"source_storage_plan"`
	DestinationStoragePlan StoragePlan `json:"DestinationStoragePlan,omitempty" bson:"destination_storage_plan"`
}
