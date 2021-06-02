package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task 存储
type TaskStorage interface {
	AddTask(t *Task) (tid primitive.ObjectID, err error)
	GetTaskList(n int) (t []*Task)
	GetTask(tid primitive.ObjectID) (*Task, error)
	SetTaskState(tid primitive.ObjectID, state TaskState) error
	SetTask(tid primitive.ObjectID, t *Task) error
	DelTask(tid primitive.ObjectID) error
	GetUserTask(uid string) (t []*Task)
	IsAllDone() bool
}
