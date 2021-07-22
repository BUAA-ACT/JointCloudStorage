package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskType = string
type TaskState string

const (
	DOWNLOAD         TaskType = "DOWNLOAD"
	DOWNLOAD_REPLICA TaskType = "DOWNLOAD_REPLICA"
	UPLOAD           TaskType = "UPLOAD"
	INDEX            TaskType = "INDEX"
	SYNC             TaskType = "SYNC"
	DELETE           TaskType = "DELETE"
	MIGRATE          TaskType = "MIGRATE"
)

const (
	CREATING   TaskState = "CREATING"   // transporter 内部状态，还不能处理
	WAITING    TaskState = "WAITING"    // 任务就绪，加入处理队列，随时可以处理
	PROCESSING TaskState = "PROCESSING" // 正在处理
	FINISH     TaskState = "FINISH"     // 成功结束
	FAIL       TaskState = "FAIL"       // 失败
	BLOCKED    TaskState = "BLOCKED"    // 阻塞，等待用户进一步的数据
)

type Task struct {
	Tid             primitive.ObjectID `bson:"_id,omitempty"`
	TaskType        TaskType           `bson:"task_type"`
	State           TaskState          `bson:"task_state"`
	StartTime       time.Time          `bson:"start_time"`
	Uid             string             `bson:"user_id"`
	SourcePath      string             `bson:"source_path"`
	DestinationPath string             `bson:"destination_path"`
	TaskOptions     *TaskOptions       `bson:"task_options"`
	Progress        float64            `bson:"progress"`
}

type TaskOptions struct {
	SourceStoragePlan      *StoragePlan `bson:"source_storage_plan"`
	DestinationStoragePlan *StoragePlan `bson:"destination_storage_plan"`
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
	if len(t.SourcePath) == 0 {
		t.SourcePath = "/"
	}
	if t.SourcePath[0] == '/' {
		return t.Uid + t.SourcePath
	} else {
		return t.Uid + "/" + t.SourcePath
	}
}
func (t *Task) GetRealDestinationPath() string {
	if len(t.DestinationPath) == 0 {
		t.DestinationPath = "/"
	}
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
		switch t.TaskOptions.DestinationStoragePlan.StorageMode {
		case StorageModelReplica:
		case StorageModelEC:
		default:
			util.Log(logrus.ErrorLevel, "task check", "wrong task type", "", string(t.TaskType), "")
			return false
		}
		if len(t.TaskOptions.DestinationStoragePlan.Clouds) == 0 {
			util.Log(logrus.ErrorLevel, "task check", "no destination cloud", "", "", "")
			return false
		}
	default:
		return true
	}
	return true
}
