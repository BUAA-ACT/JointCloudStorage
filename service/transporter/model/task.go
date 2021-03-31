package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskType string
type TaskState string

const (
	USER_UPLOAD_SIMPLE  TaskType = "USER_UPLOAD_SIMPLE"
	USER_UPLOAD_ERASURE TaskType = "USER_UPLOAD_ERASURE"
	SYNC_SIMPLE         TaskType = "SYNC_SIMPLE"
	SYNC_ERASURE        TaskType = "SYNC_ERASURE"
	DOWNLOAD            TaskType = "DOWNLOAD"
	DOWNLOAD_REPLICA    TaskType = "DOWNLOAD_REPLICA"
	UPLOAD              TaskType = "UPLOAD"
	INDEX               TaskType = "INDEX"
	SYNC                TaskType = "SYNC"
	DELETE              TaskType = "DELETE"
)

func (taskType TaskType) String() string {
	switch taskType {
	case USER_UPLOAD_SIMPLE:
		return "USER_UPLOAD_SIMPLE"
	case USER_UPLOAD_ERASURE:
		return "USER_UPLOAD_ERASURE"
	case SYNC_SIMPLE:
		return "SYNC_SIMPLE"
	case SYNC_ERASURE:
		return "SYNC_ERASURE"
	case DOWNLOAD:
		return "DOWNLOAD"
	case INDEX:
		return "INDEX"
	case UPLOAD:
		return "UPLOAD"
	case SYNC:
		return "SYNC"
	case DELETE:
		return "DELETE"
	}
	return ""
}

const (
	CREATING   TaskState = "CREATING"
	WAITING    TaskState = "WAITING"
	PROCESSING TaskState = "PROCESSING"
	FINISH     TaskState = "FINISH"
	FAIL       TaskState = "FAIL"
	BLOCKED    TaskState = "BLOCKED"
)

type Task struct {
	Tid             primitive.ObjectID `bson:"_id,omitempty"`
	TaskType        TaskType
	State           TaskState
	StartTime       time.Time
	Uid             string
	SourcePath      string
	DestinationPath string
	TaskOptions     *TaskOptions
}

type TaskOptions struct {
	SourceStoragePlan *StoragePlan
	DestinationPlan   *StoragePlan
}

type StoragePlan struct {
	StorageMode string
	Clouds      []string
	N           int
	K           int
}

func (t *Task) GetTid() primitive.ObjectID {
	return t.Tid
}

func (t *Task) GetTaskType() TaskType {
	return t.TaskType
}

func (t *Task) GetState() TaskState {
	return t.State
}

func (t *Task) GetSid() string {
	return t.Uid
}

func (t *Task) GetSourcePath() string {
	return t.SourcePath
}

func (t *Task) GetDestinationPath() string {
	return t.DestinationPath
}

func NewTask(taskType TaskType, startTime time.Time, uid string, sourcePath string, destinationPath string) *Task {

	return &Task{
		Tid:             primitive.NewObjectID(),
		TaskType:        taskType,
		State:           CREATING,
		StartTime:       startTime,
		Uid:             uid,
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
	}
}
