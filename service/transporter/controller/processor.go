package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"context"
	"errors"
	"log"
	"path"
	"reflect"
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
		case model.SYNC_SIMPLE:
			processor.taskStorage.SetTaskState(task.GetTid(), model.PROCESSING)
			go func(t *model.Task) {
				err := processor.ProcessSimpleSync(t)
				if err != nil {
					log.Panicf("Process Task Fail: %v", err)
				} else {
					log.Printf("finish task: %v", t.GetTid())
				}
				finish <- task.GetTid()
			}(task)
			log.Printf("start simple SYNC task")
		default:
			log.Fatalf("ERROR: Process TaskType: %s not implement", task.GetTaskType())
		}
	}
	for i := 0; i < len(tasks); i++ {
		<-finish
	}
}

// 处理获取临时下载 url 请求
func (processor *TaskProcessor) ProcessGetTmpDownloadUrl(t *model.Task) (url string, err error) {
	err = processor.CheckTaskType(t, model.USER_DOWNLOAD_SIMPLE)
	if err != nil {
		return "", err
	}
	storageClient := processor.storageDatabase.GetStorageClient(t.GetSid(), t.GetSourcePath())
	url, err = storageClient.GetTmpDownloadUrl(t.GetSourcePath(), time.Minute*30)
	return url, err
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
		processor.taskStorage.SetTaskState(t.GetTid(), model.FAIL)
		log.Printf("Upload user file Fail: %v", err)
		return
	}
	processor.taskStorage.SetTaskState(t.GetTid(), model.FINISH)
	return
}

// 获取用户目录信息
func (processor *TaskProcessor) ProcessPathIndex(t *model.Task) <-chan model.ObjectInfo {
	err := processor.CheckTaskType(t, model.INDEX)
	if err != nil {
		return nil
	}
	storageClient := processor.storageDatabase.GetStorageClient(t.GetSid(), t.GetSourcePath())

	return storageClient.Index(t.GetSourcePath())
}

// 普通同步任务
func (processor *TaskProcessor) ProcessSimpleSync(t *model.Task) (err error) {
	err = processor.CheckTaskType(t, model.SYNC_SIMPLE)
	if err != nil {
		return nil
	}
	// 获取源路径对应存储客户端
	storageClient := processor.storageDatabase.GetStorageClient(t.GetSid(), t.GetSourcePath())
	// 获取目的路径对应存储客户端
	destClient := processor.storageDatabase.GetStorageClient(t.GetSid(), t.GetDestinationPath())
	// 列举所有对象 todo: 如果是具有相同前缀的两个文件？
	objectCh := storageClient.Index(t.GetSourcePath())
	for obj := range objectCh {
		fileName := path.Base(obj.Key)
		// 判断两个客户端是否相同 todo: accesspoint 相同即可
		if reflect.DeepEqual(storageClient, destClient) {
			err = storageClient.Copy(obj.Key, t.GetDestinationPath()+fileName)
		} else {
			// 下载源文件
			localPath := "./tmp/" + genRandomString(10)
			err = storageClient.Download(t.GetSourcePath(), localPath)
			err = destClient.Upload(localPath, t.GetDestinationPath()+fileName)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (processor *TaskProcessor) CheckTaskType(t *model.Task, taskType model.TaskType) (err error) {
	if t.GetTaskType() != taskType {
		return errors.New("wrong task type")
	}
	return nil
}
