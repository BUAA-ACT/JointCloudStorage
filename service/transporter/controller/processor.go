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
	"strconv"
	"strings"
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
	logrus.Infof("Process %v Task Sucess, tid :%v", t.TaskType, t.Tid.Hex())
	processor.taskStorage.SetTaskState(t.Tid, model.FINISH)
}

// 处理任务
func (processor *TaskProcessor) ProcessTasks() {
	tasks := processor.taskStorage.GetTaskList(0)
	finish := make(chan primitive.ObjectID)
	for _, task := range tasks {
		processor.taskStorage.SetTaskState(task.Tid, model.PROCESSING)
		switch task.GetTaskType() {
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
				err := processor.ProcessSync(t)
				processor.SetProcessResult(t, err)
				finish <- t.Tid
			}(task)
		case model.DELETE:
			go func(t *model.Task) {
				err := processor.DeleteSingleFile(t)
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

func (processor *TaskProcessor) DeleteSingleFile(t *model.Task) error {
	fileInfo, err := processor.fileDatabase.GetFileInfo(t.Uid + "/" + t.SourcePath)
	if err != nil {
		logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		return err
	}
	var storageClients []model.StorageClient
	storageModel := t.TaskOptions.SourceStoragePlan.StorageMode
	for _, cloudName := range t.TaskOptions.SourceStoragePlan.Clouds {
		client, err := processor.storageDatabase.GetStorageClientFromName(t.Uid, cloudName)
		if err != nil {
			return err
		}
		storageClients = append(storageClients, client)
	}
	switch storageModel {
	case "Replica":
		for _, client := range storageClients {
			err = client.Remove(t.SourcePath, t.Uid)
		}
	case "EC":
		N := t.TaskOptions.SourceStoragePlan.N
		K := t.TaskOptions.SourceStoragePlan.K
		if len(storageClients) != N+K {
			return errors.New("storage num not correct")
		}
		for i, client := range storageClients {
			err = client.Remove(t.SourcePath+"."+strconv.Itoa(i), t.Uid)
		}
	default:
		return errors.New("storageModel not support")
	}
	err = processor.fileDatabase.DeleteFileInfo(fileInfo)
	return err
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
		client, err := processor.storageDatabase.GetStorageClientFromName(t.Uid, cloudName)
		if err != nil {
			return "", err
		}
		storageClients = append(storageClients, client)
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
			err := storageClients[i].Download(t.SourcePath+"."+strconv.Itoa(i), shards[i], t.Uid)
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
		err = storageClients[0].Download(t.SourcePath, rebuildPath, t.Uid)
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
	storageClient, err := processor.storageDatabase.GetStorageClientFromName(t.Uid, t.TaskOptions.SourceStoragePlan.Clouds[0])
	if err != nil {
		return "", err
	}
	url, err = storageClient.GetTmpDownloadUrl(t.GetSourcePath(), t.Uid, time.Minute*30)
	return url, err
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
			client, err := processor.storageDatabase.GetStorageClientFromName(t.Uid, cloudName)
			if err != nil {
				return err
			}
			storageClients = append(storageClients, client)
		}
		switch storageModel {
		case "Replica":
			for _, client := range storageClients {
				err = client.Upload(t.GetSourcePath(), t.GetDestinationPath(), t.Uid)
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
				err = client.Upload(shards[i], t.GetDestinationPath()+"."+strconv.Itoa(i), t.Uid)
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
	return errors.New("no storage plan")
}

// 获取用户目录信息
func (processor *TaskProcessor) ProcessPathIndex(t *model.Task) <-chan model.ObjectInfo {
	err := processor.CheckTaskType(t, model.INDEX)
	if err != nil {
		return nil
	}
	storageClient, err := processor.storageDatabase.GetStorageClientFromName(t.Uid, t.TaskOptions.SourceStoragePlan.Clouds[0])
	if err != nil {
		return nil
	}

	return storageClient.Index(t.GetSourcePath(), t.Uid)
}

// 处理同步任务
func (processor *TaskProcessor) ProcessSync(t *model.Task) (err error) {
	subTask := model.Task{}
	err = copier.Copy(&subTask, t)
	if err != nil {
		return err
	}
	// 列举所有对象
	objects, err := processor.fileDatabase.Index(t.Uid + "/" + t.SourcePath)
	if err != nil {
		return err
	}
	for _, obj := range objects {
		_, p := FromFileInfoGetUidAndPath(obj)
		subTask.SourcePath = p
		if strings.HasSuffix(t.DestinationPath, "/") {
			// 目标是一个目录
			fileName := path.Base(subTask.SourcePath)
			subTask.DestinationPath = t.DestinationPath + fileName
		}
		err = processor.ProcessSyncSingleFile(&subTask)
		if err != nil {
			return err
		}
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
	err = copier.Copy(&subTask, t)
	subTask.TaskType = model.DELETE
	subTask.TaskOptions.DestinationPlan = nil
	err = processor.DeleteSingleFile(&subTask)
	if err != nil {
		return err
	}
	logrus.Debugf("sync task %v finish", t.Tid.Hex())
	return nil
}

func (processor *TaskProcessor) CheckTaskType(t *model.Task, taskType model.TaskType) (err error) {
	if t.GetTaskType() != taskType {
		return errors.New("wrong task type")
	}
	return nil
}
