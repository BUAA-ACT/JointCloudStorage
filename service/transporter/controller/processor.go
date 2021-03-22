package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type TaskProcessor struct {
	taskStorage     model.TaskStorage
	storageDatabase model.StorageDatabase
	fileDatabase    model.FileDatabase
}

func (processor *TaskProcessor) SetTaskStorage(storage model.TaskStorage) {
	processor.taskStorage = storage
}

func (processor *TaskProcessor) SetStorageDatabase(database model.StorageDatabase) {
	processor.storageDatabase = database
}

// 创建任务
func (processor *TaskProcessor) CreateTask(taskType model.TaskType, sid string, sourcePath string, destinationPath string) {
	task := model.NewTask(taskType, time.Now(), sid, sourcePath, destinationPath)
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

func (processor *TaskProcessor) SetProcessResult(t *model.Task, err error) {
	if err != nil {
		logrus.Errorf("Process Task Fail: %v", err)
		processor.taskStorage.SetTaskState(t.Tid, model.FAIL)
	}
	processor.taskStorage.SetTaskState(t.Tid, model.FINISH)
}

// 处理任务
func (processor *TaskProcessor) ProcessTasks() {
	tasks := processor.taskStorage.GetTaskList(0)
	finish := make(chan primitive.ObjectID)
	for _, task := range tasks {
		processor.taskStorage.SetTaskState(task.Tid, model.PROCESSING)
		switch task.GetTaskType() {
		case model.USER_UPLOAD_SIMPLE:
			go func(t *model.Task) {
				err := processor.ProcessUserUploadSimple(t)
				processor.SetProcessResult(t, err)
				finish <- task.GetTid()
			}(task)
			log.Printf("start simple upload task")
		case model.SYNC_SIMPLE:
			go func(t *model.Task) {
				err := processor.ProcessSimpleSync(t)
				processor.SetProcessResult(t, err)
				finish <- t.GetTid()
			}(task)
			log.Printf("start simple SYNC task")
		case model.UPLOAD:
			go func(t *model.Task) {
				err := processor.ProcessUpload(t)
				processor.SetProcessResult(t, err)
				finish <- t.Tid
			}(task)
			log.Printf("start upload task")
		case model.DOWNLOAD:
			go func(t *model.Task) {
				filePath, err := processor.RebuildFileToDisk(t)
				if err == nil {
					err = processor.WriteDownloadUrlToDB(t, filePath)
				}
				processor.SetProcessResult(t, err)
				finish <- t.Tid
			}(task)
		case model.SYNC:
			go func(t *model.Task) {
				err := processor.ProcessSync(task)
				processor.SetProcessResult(t, err)
				finish <- t.Tid
			}(task)
		default:
			logrus.Errorf("ERROR: Process TaskType: %s not implement", task.GetTaskType())
			finish <- task.Tid
		}
	}
	for i := 0; i < len(tasks); i++ {
		id := <-finish
		logrus.Infof("finish task: %v", id.Hex())
	}
}

func (processor *TaskProcessor) WriteDownloadUrlToDB(t *model.Task, path string) error {
	fileInfo, err := processor.fileDatabase.GetFileInfo(t.Uid + "/" + t.SourcePath)
	if err != nil {
		logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		return err
	}
	accessToken, err := GenerateLocalFileAccessToken(path, t.Uid, time.Hour*24)
	if err != nil {
		logrus.Warnf("cant gen access token, err: %v", err)
	}
	fileInfo.DownloadUrl = "/cache_file?token=" + accessToken
	fileInfo.ReconstructStatus = "Done"
	fileInfo.ReconstructTime = time.Now()
	err = processor.fileDatabase.UpdateFileInfo(fileInfo)
	return err
}

func (processor *TaskProcessor) RebuildFileToDisk(t *model.Task) (path string, err error) {
	err = processor.CheckTaskType(t, model.DOWNLOAD)
	if err != nil {
		return "", err
	}
	var storageClients []model.StorageClient
	storageModel := t.TaskOptions.SourceStoragePlan.StorageMode
	for _, cloudName := range t.TaskOptions.SourceStoragePlan.Clouds {
		storageClients = append(storageClients, processor.storageDatabase.GetStorageClientFromName(cloudName, t.Uid))
	}
	if len(storageClients) == 0 {
		return "", errors.New("EC storage num wrong")
	}
	switch storageModel {
	case "EC":
		N := t.TaskOptions.SourceStoragePlan.N
		K := t.TaskOptions.SourceStoragePlan.K
		if N < 1 || K < 1 || N+K != len(storageClients) {
			return "", errors.New("EC storage num wrong")
		}
		fileInfo, err := processor.fileDatabase.GetFileInfo(t.Uid + "/" + t.SourcePath)
		if err != nil {
			logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		}
		rebuildPath := "./tmp/download/" + filepath.Base(t.SourcePath) // todo 自定义 tmp 目录
		shards := make([]string, N+K)
		for i := range shards {
			// 设置临时分块存储路径
			shards[i] = rebuildPath + fmt.Sprintf(".%d", i)
			err := storageClients[i].Download(t.SourcePath+"."+strconv.Itoa(i), shards[i])
			if err != nil {
				logrus.Errorf("Download EC block %v from %v fail: %v", shards[i], storageClients[i], err)
			}
		}
		err = Decode(rebuildPath, fileInfo.Size, shards, N, K) // todo 文件大小
		if err != nil {
			logrus.Errorf("Rebuild File %v fail: %v", rebuildPath, err)
			return "", err
		}
		return rebuildPath, nil
	case "Replica":
		_, err := processor.fileDatabase.GetFileInfo(t.Uid + "/" + t.SourcePath)
		if err != nil {
			logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		}
		rebuildPath := "./tmp/download/" + filepath.Base(t.SourcePath) // todo 自定义 tmp 目录
		err = storageClients[0].Download(t.SourcePath, rebuildPath)
		if err != nil {
			logrus.Errorf("Download Replica %v from %v fail: %v", t.SourcePath, storageClients[0], err)
		}
		return rebuildPath, nil
	default:
		return "", errors.New("storageModel not support")
	}
}

// 处理获取临时下载 url 请求
func (processor *TaskProcessor) ProcessGetTmpDownloadUrl(t *model.Task) (url string, err error) {
	err = processor.CheckTaskType(t, model.DOWNLOAD_REPLICA)
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

func (processor *TaskProcessor) ProcessUpload(t *model.Task) (err error) {
	if t.GetTaskType() != model.UPLOAD {
		return errors.New("wrong task type")
	}
	if t.GetState() == model.FINISH {
		return errors.New("task already finish")
	}
	// 判断上传方式
	var storageClients []model.StorageClient
	if t.TaskOptions != nil {
		storageModel := t.TaskOptions.DestinationPlan.StorageMode
		for _, cloudName := range t.TaskOptions.DestinationPlan.Clouds {
			storageClients = append(storageClients, processor.storageDatabase.GetStorageClientFromName(cloudName, t.Uid))
		}
		switch storageModel {
		case "Replica":
			for _, client := range storageClients {
				err = client.Upload(t.GetSourcePath(), t.GetDestinationPath())
			}
			fileInfo, err := model.NewFileInfoFromPath(t.SourcePath, t.Uid, t.DestinationPath)
			if CheckErr(err, "New File Info") {
				return err
			}
			err = processor.fileDatabase.CreateFileInfo(fileInfo)
			CheckErr(err, "Create File Info")
		case "EC": // 纠删码模式
			N := t.TaskOptions.DestinationPlan.N
			K := t.TaskOptions.DestinationPlan.K
			if N < 1 || K < 1 || N+K != len(storageClients) {
				return errors.New("EC storage num wrong")
			}
			shards := make([]string, N+K)
			for i := range shards {
				// 设置临时分块存储路径
				shards[i] = t.GetSourcePath() + fmt.Sprintf(".%d", i)
			}
			// 开始分块
			err := Encode(t.GetSourcePath(), shards, N, K)
			if err != nil {
				logrus.Errorf("Encode file %s failed.", t.GetSourcePath())
				return err
			}
			// 开始上传
			for i, client := range storageClients {
				err = client.Upload(shards[i], t.GetDestinationPath()+"."+strconv.Itoa(i))
			}
			if CheckErr(err, "Upload EC block") {
				return err
			}
			fileInfo, err := model.NewFileInfoFromPath(t.SourcePath, t.Uid, t.DestinationPath)
			if CheckErr(err, "New File Info") {
				return err
			}
			err = processor.fileDatabase.CreateFileInfo(fileInfo)
			CheckErr(err, "Create File Info")
		default:
			return errors.New("storage model not implement")
		}
		return
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

// 处理同步任务
func (processor *TaskProcessor) ProcessSync(t *model.Task) (err error) {
	subTask := model.Task{}
	err = copier.Copy(&subTask, t)
	if err != nil {
		return err
	}
	return nil
}

// 处理同步任务
func (processor *TaskProcessor) ProcessSyncSingleFile(t *model.Task) (err error) {
	// 从源端下载文件到本地
	subTask := model.Task{}
	err = copier.Copy(&subTask, t)
	if err != nil {
		return err
	}
	subTask.TaskType = model.DOWNLOAD
	filePath, err := processor.RebuildFileToDisk(&subTask)
	if err != nil {
		return err
	}
	logrus.Debugf("rebuile file finish, path: %v", filePath)
	subTask.SourcePath = filePath
	subTask.TaskOptions.SourceStoragePlan = nil
	subTask.TaskType = model.UPLOAD
	err = processor.ProcessUpload(&subTask)
	if err != nil {
		return err
	}
	logrus.Debugf("sync task %v finish", t.Tid.Hex())
	return nil
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
