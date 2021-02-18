package transporter

import (
	"errors"
	"log"
	"time"
)

type TaskProcessor struct {
	taskStorage     TaskStorage
	storageDatabase StorageDatabase
}

func (processor *TaskProcessor) SetTaskStorage(storage TaskStorage) {
	processor.taskStorage = storage
}

func (processor *TaskProcessor) SetStorageDatabase(database StorageDatabase) {
	processor.storageDatabase = database
}

// 创建任务
func (processor *TaskProcessor) CreateTask(taskType TaskType, sid string, sourcePath string, destinationPath string) {
	task := Task{tid: 0, taskType: taskType, state: CREATING, startTime: time.Now(), sid: sid, sourcePath: sourcePath, destinationPath: destinationPath}
	_, err := processor.taskStorage.AddTask(task)
	if err != nil {
		log.Panicf("Create Task ERROR: %v", err)
	}
}

// 处理任务
func (processor *TaskProcessor) ProcessTasks() {
	tasks := processor.taskStorage.GetTaskList(0)
	for _, task := range tasks {
		switch task.taskType {
		case USER_UPLOAD_SIMPLE:
			err := processor.ProcessUserUploadSimple(task)
			if err != nil {
				log.Panicf("Process Task Fail: %v", err)
			}
			log.Printf("start simple upload task")
		default:
			log.Fatalf("ERROR: Process TaskType: %s not implement", task.taskType)
		}
	}
}

func (processor *TaskProcessor) ProcessUserUploadSimple(t Task) (err error) {
	if t.taskType != USER_UPLOAD_SIMPLE {
		return errors.New("wrong task type")
	}
	if t.state == FINISH {
		return errors.New("task already finish")
	}
	storageClient := processor.storageDatabase.GetStorageClient(t.sid, t.destinationPath)
	err = storageClient.Upload(t.sourcePath, t.destinationPath)
	return
}
