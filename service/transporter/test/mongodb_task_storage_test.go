package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestMongoTaskStorage_AddTask(t *testing.T) {
	mongots, err := model.NewMongoTaskStorage()
	if err != nil {
		t.Error("can't connect to the mongodb")
	}
	var tid primitive.ObjectID
	task := model.NewTask(model.UPLOAD, time.Now(), "asdf", "asdfasdf", "asdfasdf")
	task.TaskOptions = &model.TaskOptions{
		SourceStoragePlan: &model.StoragePlan{
			StorageMode: "asdfasd",
			Clouds:      []string{"asdfasd", "asdfasdfa"},
		},
		DestinationStoragePlan: &model.StoragePlan{
			StorageMode: "asdfasdfasdf",
			Clouds:      nil,
		},
	}
	t.Run("test the AddTask", func(t *testing.T) {
		tid, err = mongots.AddTask(task)
		if err != nil {
			err.Error()
		}
		fmt.Println(tid)
	})
	t.Run("test the gettask", func(t *testing.T) {
		task, err := mongots.GetTask(tid)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(*task)
	})
	t.Run("test the setState", func(t *testing.T) {
		err = mongots.SetTaskState(tid, model.WAITING)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("test the gettasklist", func(t *testing.T) {
		res := mongots.GetTaskList(3)
		for _, re := range res {
			fmt.Println(*re)
		}
	})

	t.Run("test the deleteTask", func(t *testing.T) {
		err = mongots.DelTask(tid)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("test the isAllDone", func(t *testing.T) {
		a := mongots.IsAllDone()
		if a {
			fmt.Println("not all done")
		}
	})
}
