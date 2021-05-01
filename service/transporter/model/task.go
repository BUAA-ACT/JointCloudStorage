package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"github.com/sirupsen/logrus"
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
	UPLOAD              TaskType = "Upload"
	INDEX               TaskType = "Index"
	SYNC                TaskType = "Sync"
	DELETE              TaskType = "Delete"
	MIGRATE             TaskType = "Migrate"
)

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

type StorageModel string

const (
	StorageModelReplica StorageModel = "Replica"
	StorageModelEC      StorageModel = "EC"
	StorageModelMigrate StorageModel = "Migrate"
)

type StoragePlan struct {
	StorageMode StorageModel
	Clouds      []string
	N           int
	K           int
}

func (t *Task) GetRealSourcePath() string {
	if t.SourcePath[0] == '/' {
		return t.Uid + t.SourcePath
	} else {
		return t.Uid + "/" + t.SourcePath
	}
}
func (t *Task) GetRealDestinationPath() string {
	if t.DestinationPath[0] == '/' {
		return t.Uid + t.DestinationPath
	} else {
		return t.Uid + "/" + t.DestinationPath
	}
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

func (t *Task) Check() bool {
	logrus.Tracef("%+v", t)
	logrus.Tracef("%+v", t.TaskOptions)
	switch t.TaskType {
	case UPLOAD:
		if t.DestinationPath == "" {
			return false
		}
		switch t.TaskOptions.DestinationPlan.StorageMode {
		case StorageModelReplica:
		case StorageModelEC:
		default:
			util.Log(logrus.ErrorLevel, "task check", "wrong task type", "", string(t.TaskType), "")
			return false
		}
		if len(t.TaskOptions.DestinationPlan.Clouds) == 0 {
			util.Log(logrus.ErrorLevel, "task check", "no destination cloud", "", "", "")
			return false
		}
	default:
		return true
	}
	return true
}
