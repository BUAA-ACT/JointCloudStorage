package model

import (
	"errors"
	"github.com/jinzhu/copier"
	"sync"
	"time"
)

type TaskType int32
type TaskState int32

const (
	USER_UPLOAD_SIMPLE   TaskType = 1
	USER_UPLOAD_ERASURE  TaskType = 2
	SYNC_SIMPLE          TaskType = 3
	SYNC_ERASURE         TaskType = 4
	USER_DOWNLOAD_SIMPLE TaskType = 5
	UPLOAD               TaskType = 6
	INDEX                TaskType = 7
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
	case USER_DOWNLOAD_SIMPLE:
		return "USER_DOWNLOAD_SIMPLE"
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
	Tid             int
	TaskType        TaskType
	State           TaskState
	StartTime       time.Time
	Sid             string
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
}

func (t *Task) GetTid() int {
	return t.Tid
}

func (t *Task) GetTaskType() TaskType {
	return t.TaskType
}

func (t *Task) GetState() TaskState {
	return t.State
}

func (t *Task) GetSid() string {
	return t.Sid
}

func (t *Task) GetSourcePath() string {
	return t.SourcePath
}

func (t *Task) GetDestinationPath() string {
	return t.DestinationPath
}

func NewTask(tid int, taskType TaskType, startTime time.Time, sid string, sourcePath string, destinationPath string) *Task {
	return &Task{
		Tid:             0,
		TaskType:        taskType,
		State:           CREATING,
		StartTime:       startTime,
		Sid:             sid,
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
	}
}

// Task 存储
type TaskStorage interface {
	AddTask(t *Task) (tid int, err error)
	GetTaskList(n int) (t []*Task)
	GetTask(tid int) (t *Task, err error)
	SetTaskState(tid int, state TaskState) (err error)
	SetTask(tid int, t *Task) (err error)
	DelTask(tid int) (err error)
}

type InMemoryTaskStorage struct {
	taskList []*Task
	maxTid   int
	mutex    sync.Mutex
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	return &InMemoryTaskStorage{
		taskList: make([]*Task, 0),
		maxTid:   0,
		mutex:    sync.Mutex{},
	}
}

func (s *InMemoryTaskStorage) GetTask(tid int) (t *Task, err error) {
	for _, task := range s.taskList {
		if task.Tid == tid {
			return task, nil
		}
	}
	return nil, errors.New("task not found")
}

func (s *InMemoryTaskStorage) SetTask(tid int, t *Task) (err error) {
	for i, task := range s.taskList {
		if task.Tid == tid {
			newTask := Task{}
			copier.Copy(&newTask, t)
			s.taskList[i] = &newTask
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *InMemoryTaskStorage) AddTask(t *Task) (tid int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t.Tid = s.maxTid + 1
	s.maxTid += 1
	if t.State != BLOCKED {
		t.State = WAITING
	}
	s.taskList = append(s.taskList, t)
	return t.Tid, nil
}

func (s *InMemoryTaskStorage) GetTaskList(n int) (t []*Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, task := range s.taskList {
		if task.State == WAITING {
			t = append(t, task)
		}
	}
	return
}
func (s *InMemoryTaskStorage) SetTaskState(tid int, state TaskState) (err error) {
	for i, task := range s.taskList {
		if task.Tid == tid {
			s.taskList[i].State = state
			return nil
		}
	}
	return errors.New("task not found")
}
func (s *InMemoryTaskStorage) DelTask(tid int) (err error) {
	for i, task := range s.taskList {
		if task.Tid == tid {
			s.taskList = append(s.taskList[:i], s.taskList[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
