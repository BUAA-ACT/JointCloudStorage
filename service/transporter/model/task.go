package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskType int32
type TaskState int32

const (
	USER_UPLOAD_SIMPLE  TaskType = 1
	USER_UPLOAD_ERASURE TaskType = 2
	SYNC_SIMPLE         TaskType = 3
	SYNC_ERASURE        TaskType = 4
	DOWNLOAD_EC         TaskType = 5
	DOWNLOAD_REPLICA    TaskType = 6
	UPLOAD              TaskType = 7
	INDEX               TaskType = 8
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
	case DOWNLOAD_EC:
		return "DOWNLOAD_EC"
	case INDEX:
		return "INDEX"
	}
	return ""
}

const (
	CREATING   TaskState = 1
	WAITING    TaskState = 2
	PROCESSING TaskState = 3
	FINISH     TaskState = 4
	FAIL       TaskState = 5
	BLOCKED    TaskState = 6
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
