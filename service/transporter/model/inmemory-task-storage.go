package model

import (
	"errors"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

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

func (s *InMemoryTaskStorage) GetTask(tid primitive.ObjectID) (t *Task, err error) {
	for _, task := range s.taskList {
		if task.Tid == tid {
			return task, nil
		}
	}
	return nil, errors.New("task not found")
}

func (s *InMemoryTaskStorage) SetTask(tid primitive.ObjectID, t *Task) (err error) {
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

func (s *InMemoryTaskStorage) AddTask(t *Task) (tid primitive.ObjectID, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t.Tid = primitive.NewObjectID()
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
func (s *InMemoryTaskStorage) SetTaskState(tid primitive.ObjectID, state TaskState) (err error) {
	for i, task := range s.taskList {
		if task.Tid == tid {
			s.taskList[i].State = state
			return nil
		}
	}
	return errors.New("task not found")
}
func (s *InMemoryTaskStorage) DelTask(tid primitive.ObjectID) (err error) {
	for i, task := range s.taskList {
		if task.Tid == tid {
			s.taskList = append(s.taskList[:i], s.taskList[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *InMemoryTaskStorage) IsAllDone() bool {
	for _, task := range s.taskList {
		if task.State != FINISH && task.State != FAIL {
			return false
		}
	}
	return true
}

func (s *InMemoryTaskStorage) GetUserTask(uid string) (t []*Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, task := range s.taskList {
		if task.Uid == uid {
			t = append(t, task)
		}
	}
	return
}
