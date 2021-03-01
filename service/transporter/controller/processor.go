package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"context"
	"errors"
	"log"
	"time"
)

type TaskProcessor struct {
	taskStorage     model.TaskStorage
	storageDatabase model.StorageDatabase
}

func (processor *TaskProcessor) SetTaskStorage(storage model.TaskStorage) {
	processor.taskStorage = storage
}

func (processor *TaskProcessor) SetStorageDatabase(database model.StorageDatabase) {
	processor.storageDatabase = database
}

// 创建任务
func (processor *TaskProcessor) CreateTask(taskType model.TaskType, sid string, sourcePath string, destinationPath string) {
	task := model.NewTask(0, taskType, time.Now(), sid, sourcePath, destinationPath)
	_, err := processor.taskStorage.AddTask(task)
	if err != nil {
		log.Panicf("Create Task ERROR: %v", err)
	}
}

func (processor *TaskProcessor) StartProcessTasks(ctx context.Context) {
	go func() {
		for true {
			select {
			case <-ctx.Done():
				log.Fatal("Processor stop")
				return
			default:
			}
			go processor.ProcessTasks()
			time.Sleep(time.Second)
		}
	}()
}

// 处理任务
func (processor *TaskProcessor) ProcessTasks() {
	tasks := processor.taskStorage.GetTaskList(0)
	var finish chan int
	for _, task := range tasks {
		switch task.GetTaskType() {
		case model.USER_UPLOAD_SIMPLE:
			processor.taskStorage.SetTaskState(task.GetTid(), model.PROCESSING)
			go func(t *model.Task) {
				err := processor.ProcessUserUploadSimple(t)
				if err != nil {
					log.Panicf("Process Task Fail: %v", err)
				} else {
					log.Printf("finish task: %v", t.GetTid())
				}
				finish <- task.GetTid()
			}(task)
			log.Printf("start simple upload task")
		default:
			log.Fatalf("ERROR: Process TaskType: %s not implement", task.GetTaskType())
		}
	}
	for i := 0; i < len(tasks); i++ {
		<-finish
	}
}

// 处理用户上传
func (processor *TaskProcessor) ProcessUserUploadSimple(t *model.Task) (err error) {
	if t.GetTaskType() != model.USER_UPLOAD_SIMPLE {
		return errors.New("wrong task type")
	}
	if t.GetState() == model.FINISH {
		return errors.New("task already finish")
	}
	// 先获取当前用户上传路径对应的存储客户端
	storageClient := processor.storageDatabase.GetStorageClient(t.GetSid(), t.GetDestinationPath())
	// 存储客户端上传文件
	err = storageClient.Upload(t.GetSourcePath(), t.GetDestinationPath())
	if err != nil {
		log.Panicf("Upload user file Fail: %v", err)
	}
	processor.taskStorage.SetTaskState(t.GetTid(), model.FINISH)
	return
}
