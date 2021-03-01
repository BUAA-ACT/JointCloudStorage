package model

import (
	"errors"
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
	INDEX                TaskType = 6
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
)

type Task struct {
	tid             int
	taskType        TaskType
	state           TaskState
	startTime       time.Time
	sid             string
	sourcePath      string
	destinationPath string
}

func (t *Task) GetTid() int {
	return t.tid
}

func (t *Task) GetTaskType() TaskType {
	return t.taskType
}

func (t *Task) GetState() TaskState {
	return t.state
}

func (t *Task) GetSid() string {
	return t.sid
}

func (t *Task) GetSourcePath() string {
	return t.sourcePath
}

func (t *Task) GetDestinationPath() string {
	return t.destinationPath
}

func NewTask(tid int, taskType TaskType, startTime time.Time, sid string, sourcePath string, destinationPath string) *Task {
	return &Task{
		tid:             0,
		taskType:        taskType,
		state:           CREATING,
		startTime:       startTime,
		sid:             sid,
		sourcePath:      sourcePath,
		destinationPath: destinationPath,
	}
}

// Task 存储
type TaskStorage interface {
	AddTask(t *Task) (tid int, err error)
	GetTaskList(n int) (t []*Task)
	SetTaskState(tid int, state TaskState) (err error)
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

func (s *InMemoryTaskStorage) AddTask(t *Task) (tid int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t.tid = s.maxTid + 1
	s.maxTid += 1
	t.state = WAITING
	s.taskList = append(s.taskList, t)
	return t.tid, nil
}

func (s *InMemoryTaskStorage) GetTaskList(n int) (t []*Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, task := range s.taskList {
		if task.state == WAITING {
			t = append(t, task)
		}
	}
	return
}
func (s *InMemoryTaskStorage) SetTaskState(tid int, state TaskState) (err error) {
	for i, task := range s.taskList {
		if task.tid == tid {
			s.taskList[i].state = state
			return nil
		}
	}
	return errors.New("task not found")
}
func (s *InMemoryTaskStorage) DelTask(tid int) (err error) {
	for i, task := range s.taskList {
		if task.tid == tid {
			s.taskList = append(s.taskList[:i], s.taskList[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
