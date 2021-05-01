package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

type TaskProcessor struct {
	taskStorage   model.TaskStorage
	cloudDatabase model.CloudDatabase
	FileDatabase  model.FileDatabase
	Lock          *Lock
	Scheduler     Scheduler
	Monitor       *TrafficMonitor
	UserDatabase  model.UserDatabase
}

func (processor *TaskProcessor) SetTaskStorage(storage model.TaskStorage) {
	processor.taskStorage = storage
}

func (processor *TaskProcessor) SetStorageDatabase(database model.CloudDatabase) {
	processor.cloudDatabase = database
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
	} else {
		logrus.Infof("Process %v Task Sucess, tid :%v", t.TaskType, t.Tid.Hex())
		processor.taskStorage.SetTaskState(t.Tid, model.FINISH)
	}
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
					err = processor.WriteDownloadUrlToDB(t, filePath, t.TaskOptions.SourceStoragePlan.Clouds[0])
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
				err := processor.DeleteStorageFile(t)
				processor.SetProcessResult(t, err)
				finish <- t.Tid
			}(task)
		case model.MIGRATE:
			go func(t *model.Task) {
				err := processor.ProcessMigrate(t)
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

// DeleteFileInfo 删除文件元信息
func (processor *TaskProcessor) DeleteFileInfo(t *model.Task) error {
	fileInfo, err := processor.FileDatabase.GetFileInfo(t.GetRealSourcePath())
	if err != nil {
		logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		return err
	}
	err = processor.FileDatabase.DeleteFileInfo(fileInfo)
	return err
}

// DeleteStorageFile 删除文件存储
func (processor *TaskProcessor) DeleteStorageFile(t *model.Task) (err error) {
	var storageClients []model.StorageClient
	storageModel := t.TaskOptions.SourceStoragePlan.StorageMode
	for _, cloudName := range t.TaskOptions.SourceStoragePlan.Clouds {
		client, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, cloudName)
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
	return
}

// 弃用
func (processor *TaskProcessor) DeleteSingleFile(t *model.Task) error {
	fileInfo, err := processor.FileDatabase.GetFileInfo(t.GetRealSourcePath())
	if err != nil {
		logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		return err
	}
	var storageClients []model.StorageClient
	storageModel := t.TaskOptions.SourceStoragePlan.StorageMode
	for _, cloudName := range t.TaskOptions.SourceStoragePlan.Clouds {
		client, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, cloudName)
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
	err = processor.FileDatabase.DeleteFileInfo(fileInfo)
	return err
}

func (processor *TaskProcessor) WriteDownloadUrlToDB(t *model.Task, localFilePath string, cloudID string) error {
	fileInfo, err := processor.FileDatabase.GetFileInfo(t.GetRealSourcePath())
	if err != nil {
		logrus.Warnf("cant get file info: %v%v, err: %v", t.Uid, t.SourcePath, err)
		return err
	}
	fileName := path.Base(fileInfo.FileID)
	accessToken, err := util.GenerateLocalFileAccessToken(localFilePath, t.Uid, time.Hour*24, fileName)
	if err != nil {
		logrus.Warnf("cant gen access token, err: %v", err)
	}
	fileInfo.DownloadUrl = "/cache_file?token=" + accessToken
	fileInfo.ReconstructStatus = model.FileDone
	fileInfo.LastReconstructed = time.Now()
	err = processor.FileDatabase.UpdateFileInfo(fileInfo)
	_, err = processor.Monitor.AddDownloadTraffic(t.Uid, fileInfo.Size, cloudID)
	if err != nil {
		return err
	}
	return err
}

func (processor *TaskProcessor) RebuildFileToDisk(t *model.Task) (path string, err error) {
	err = processor.CheckTaskType(t, model.DOWNLOAD)
	if err != nil {
		return "", err
	}
	fileInfo, err := processor.FileDatabase.GetFileInfo(t.GetRealSourcePath())
	if err != nil {
		util.Log(logrus.ErrorLevel, "processor", "rebuild can't get fileInfo", t.GetRealSourcePath(), "", err.Error())
		return "", errors.New(util.ErrorMsgCantGetFileInfo)
	}
	fileInfo.ReconstructStatus = model.FileWorking
	_ = processor.FileDatabase.UpdateFileInfo(fileInfo)
	var storageClients []model.StorageClient
	storageModel := t.TaskOptions.SourceStoragePlan.StorageMode
	for _, cloudName := range t.TaskOptions.SourceStoragePlan.Clouds {
		client, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, cloudName)
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
		rebuildPath := util.Config.DownloadFileTempPath + util.GenRandomString(20)
		shards := make([]string, N+K)
		for i := range shards {
			// 设置临时分块存储路径
			shards[i] = rebuildPath + fmt.Sprintf(".%d", i)
			err := storageClients[i].Download(t.SourcePath+"."+strconv.Itoa(i), shards[i], t.Uid)
			if err != nil {
				logrus.Errorf("Download EC block %v from %v fail: %v", shards[i], storageClients[i], err)
				shards[i] = shards[i] + ".fail"
				continue
			}
			processor.Monitor.AddUploadTrafficFromFile(t.Uid, shards[i], t.TaskOptions.SourceStoragePlan.Clouds[i])
		}
		err = Decode(rebuildPath, fileInfo.Size, shards, N, K)
		if err != nil {
			logrus.Errorf("Rebuild File %v fail: %v", rebuildPath, err)
			return "", err
		}
		return rebuildPath, nil
	case "Replica":
		rebuildPath := util.Config.DownloadFileTempPath + util.GenRandomString(20)
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
	storageClient, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, t.TaskOptions.SourceStoragePlan.Clouds[0])
	if err != nil {
		return "", err
	}
	url, err = storageClient.GetTmpDownloadUrl(t.GetSourcePath(), t.Uid, time.Minute*30)
	if err != nil {
		return "", err
	}
	return url, err
}

// AddFileInfo 增加 fileInfo 到数据库
func (processor *TaskProcessor) AddFileInfo(t *model.Task) (err error) {
	_, fileInfoErr := processor.FileDatabase.GetFileInfo(t.GetRealDestinationPath())
	if fileInfoErr == nil {
		return nil
	}
	fileInfo, err := model.NewFileInfoFromPath(t.SourcePath, t.Uid, t.DestinationPath)
	fileInfo.SyncStatus = model.FilePending
	if err != nil {
		return err
	}
	err = processor.FileDatabase.CreateFileInfo(fileInfo)
	return err
}

func (processor *TaskProcessor) ProcessUpload(t *model.Task) (err error) {
	if t.GetTaskType() != model.UPLOAD {
		return errors.New("wrong task type")
	}
	if t.GetState() == model.FINISH {
		return errors.New("task already finish")
	}
	defer processor.Lock.UnLock(t.GetRealDestinationPath())
	fileInfo, fileInfoErr := processor.FileDatabase.GetFileInfo(t.GetRealDestinationPath())
	if fileInfoErr == nil { // 更新文件同步状态
		fileInfo.SyncStatus = model.FileWorking
		_ = processor.FileDatabase.UpdateFileInfo(fileInfo)
	}
	// 判断上传方式
	var storageClients []model.StorageClient
	if t.TaskOptions != nil {
		storageModel := t.TaskOptions.DestinationPlan.StorageMode
		for _, cloudName := range t.TaskOptions.DestinationPlan.Clouds {
			client, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, cloudName)
			if err != nil {
				return err
			}
			storageClients = append(storageClients, client)
		}
		switch storageModel {
		case "Replica":
			fileInfo, err = model.NewFileInfoFromPath(t.SourcePath, t.Uid, t.DestinationPath)
			if util.CheckErr(err, "New File Info") {
				return err
			}
			for i, client := range storageClients {
				_, err = processor.Monitor.AddUploadTraffic(t.Uid, fileInfo.Size, t.TaskOptions.DestinationPlan.Clouds[i])
				err = client.Upload(t.GetSourcePath(), t.GetDestinationPath(), t.Uid)
			}
		case "EC": // 纠删码模式
			fileInfo, err = model.NewFileInfoFromPath(t.SourcePath, t.Uid, t.DestinationPath)
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
			err = Encode(t.GetSourcePath(), shards, N, K)
			if err != nil {
				logrus.Errorf("Encode file %s failed.", t.GetSourcePath())
				return err
			}
			// 开始上传
			for i, client := range storageClients {
				err = client.Upload(shards[i], t.GetDestinationPath()+"."+strconv.Itoa(i), t.Uid)
				if err != nil {
					util.Log(logrus.ErrorLevel, "process upload EC",
						"client upload fail", "", "", err.Error())
					continue
				}
				processor.Monitor.AddUploadTrafficFromFile(t.Uid, shards[i], t.TaskOptions.DestinationPlan.Clouds[i])
			}
		default:
			return errors.New("storage model not implement")
		}
		// 上传后，更新 Sync Status， 更新流量统计
		if util.CheckErr(err, "Upload file to cloud") {
			return err
		}
		fileInfo.LastModified = time.Now()
		if err != nil {
			fileInfo.SyncStatus = model.FileFail
		} else {
			fileInfo.SyncStatus = model.FileDone
		}
		if fileInfoErr != nil { // 文件之前不存在
			err = processor.FileDatabase.CreateFileInfo(fileInfo)
		} else {
			err = processor.FileDatabase.UpdateFileInfo(fileInfo)
		}
		if util.CheckErr(err, "Create File Info") {
			return err
		}
		_, err = processor.Monitor.AddVolume(t.Uid, fileInfo.Size)
		err := processor.Scheduler.UploadFileMetadata(t.TaskOptions.DestinationPlan.Clouds, t.Uid, fileInfo) // todo 此处错误被隐藏
		util.CheckErr(err, "File Metadata sync")
	} else {
		return errors.New("no storage plan")
	}
	return err
}

// 获取用户目录信息
func (processor *TaskProcessor) ProcessPathIndex(t *model.Task) <-chan model.ObjectInfo {
	err := processor.CheckTaskType(t, model.INDEX)
	if err != nil {
		return nil
	}
	storageClient, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, t.TaskOptions.SourceStoragePlan.Clouds[0])
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
	objects, err := processor.FileDatabase.Index(t.GetRealSourcePath())
	if err != nil {
		return err
	}
	for _, obj := range objects {
		_, p := model.FromFileInfoGetUidAndPath(obj)
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
	user, err := processor.UserDatabase.GetUserFromID(t.Uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "processor", "can't get userInfo when process sync", t.Uid, "", err.Error())
		return err
	}
	user.Status = model.NormalUser
	err = processor.UserDatabase.UpdateUserInfo(user)
	return err
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
	logrus.Debugf("rebuild file finish, path: %v", filePath)
	subTask.SourcePath = filePath
	subTask.TaskOptions.SourceStoragePlan = nil
	subTask.TaskType = model.UPLOAD
	err = processor.Lock.Lock(subTask.GetRealDestinationPath()) // todo 在这里加锁是否是有必要的
	if err != nil {
		logrus.Errorf("get file Lock fail: %v", err)
	}
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

// 处理简单迁移任务
func (processor *TaskProcessor) ProcessMigrate(t *model.Task) (err error) {
	if len(t.TaskOptions.SourceStoragePlan.Clouds) != len(t.TaskOptions.DestinationPlan.Clouds) {
		return errors.New(util.ErrorMsgWrongCloudNum)
	}
	for i, sourceCloudID := range t.TaskOptions.SourceStoragePlan.Clouds {
		destCloudID := t.TaskOptions.DestinationPlan.Clouds[i]
		srcClient, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, sourceCloudID)
		dstClient, err := processor.cloudDatabase.GetStorageClientFromName(t.Uid, destCloudID)
		if err != nil {
			return err
		}
		objectsChan := srcClient.Index(t.SourcePath, t.Uid)
		for object := range objectsChan {
			rebuildPath := util.Config.DownloadFileTempPath + util.GenRandomString(20)
			err = srcClient.Download(object.Key, rebuildPath, t.Uid)
			if err != nil {
				logrus.Errorf("Download Replica %v from %v fail: %v", t.SourcePath, srcClient, err)
				return errors.New(util.ErrorMsgProcessMigrateDownloadErr)
			}
			err = dstClient.Upload(rebuildPath, object.Key, t.Uid)
			if err != nil {
				logrus.Errorf("Upload Replica %v from %v fail: %v", t.SourcePath, srcClient, err)
				return errors.New(util.ErrorMsgProcessMigrateUploadErr)
			}
		}
	}
	return nil //todo
}

func (processor *TaskProcessor) CheckTaskType(t *model.Task, taskType model.TaskType) (err error) {
	if t.GetTaskType() != taskType {
		return errors.New("wrong task type")
	}
	return nil
}
