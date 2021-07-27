package model

import "time"

type Task struct {
	TaskID        string    `json:"TaskID" bson:"_id,omitempty"`
	TaskType      string    `json:"TaskType" bson:"task_type"`
	TaskState     string    `json:"TaskState" bson:"task_state"`
	TaskStartTime time.Time `json:"TaskStartTime" bson:"start_time"`
	UserID        string    `json:"UserID" bson:"user_id"`
	Progress      float64   `json:"Progress" bson:"progress"`
}
