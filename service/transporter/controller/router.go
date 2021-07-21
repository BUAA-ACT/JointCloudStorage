package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Router struct {
	*gin.Engine
	processor TaskProcessor
}

type RequestTask struct {
	TaskType               model.TaskType        `json:"TaskType"`
	Uid                    string                `json:"UserID"`
	DestinationPath        string                `json:"DestinationPath"`
	SourcePath             string                `json:"SourcePath"`
	SourceStoragePlan      model.UserStoragePlan `json:"SourceStoragePlan"`
	DestinationStoragePlan model.UserStoragePlan `json:"DestinationStoragePlan"`
}

type RequestCloud struct {
	ID string `json:"ID"`
}

type RequestTaskReply struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data TaskResult `json:"data"`
}

type TaskResult struct {
	Type   string
	Result string
}

type RequestGetUserTaskReply struct {
	Code  int
	Msg   string
	Tasks []*model.Task
}

func NewRouter(processor TaskProcessor) *Router {
	var router Router
	engine := gin.Default()
	engine.Use(util.CORSMiddleware())
	router = Router{
		Engine:    engine,
		processor: processor,
	}
	router.GET("/", Index)
	router.POST("/upload/*path", util.JWTAuthMiddleware(), router.AddUploadTask)
	router.POST("/upload", util.JWTAuthMiddleware(), router.AddUploadTask)
	router.GET("/jcspan/*path", router.GetFile)
	router.GET("/index/*path", router.FileIndex)
	router.POST("/task", router.CreateTask)
	router.GET("/cache_file", util.JWTAuthMiddleware(), router.GetLocalFileByToken)
	router.GET("/state/:key", router.GetState)
	router.GET("/debug/:key", router.Debug)
	router.GET("/task/:uid", router.GetUserTask)
	rand.Seed(time.Now().Unix())
	return &router
}

func (router *Router) crossDomain(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT")
}

func (router *Router) Debug(c *gin.Context) {
	key := c.Param("key")
	switch key {
	case "unlock_test_user":
		router.processor.Lock.UnLockAll("/tester")
	case "drop_task_table":
		util.ClearAll()
	case "get_file_download_url":
		fileID := c.Query("id")
		info, err := router.processor.FileDatabase.GetFileInfo(fileID)
		if err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}
		c.String(http.StatusOK, info.DownloadUrl)
	default:
		c.String(http.StatusBadRequest, "key not imply")
	}
}

func (router *Router) GetState(c *gin.Context) {
	key := c.Param("key")
	switch key {
	case "process_state":
		if router.processor.TaskStorage.IsAllDone() {
			c.String(http.StatusOK, "done")
		} else {
			c.String(http.StatusOK, "working")
		}
	default:
		c.String(http.StatusBadRequest, "key not imply")
	}
}

func (router *Router) GetLocalFileByToken(c *gin.Context) {
	filePath := c.MustGet("filePath").(string)
	fileName, exist := c.Get("fileName")
	var fileNameStr string
	if !exist {
		fileNameStr = "File"
	} else {
		fileNameStr = fileName.(string)
	}
	util.Printf("Router", "receive get local file req, file name: %v", fileNameStr)
	mime, _ := mimetype.DetectFile(filePath)
	c.Header("Content-Type", mime.String())
	c.Header("Content-Disposition", "attachment;filename="+fileNameStr)
	c.File(filePath)
}

func (router *Router) CreateTask(c *gin.Context) {
	var reqTask RequestTask

	err := json.NewDecoder(c.Request.Body).Decode(&reqTask)
	if err != nil {
		reply := RequestTaskReply{
			Code: util.ErrorCodeWrongRequestFormat,
			Msg:  util.ErrorMsgWrongRequestFormat,
			Data: TaskResult{},
		}
		c.JSON(util.ErrorCodeWrongRequestFormat, reply)
		return
	}
	reqTask.TaskType = strings.ToUpper(reqTask.TaskType)

	switch reqTask.TaskType {
	case model.UPLOAD:
		task := RequestTask2Task(&reqTask, model.UPLOAD, model.BLOCKED)
		if !task.Check() {
			taskRequestReplyErr(util.ErrorCodeWrongRequestFormat, util.ErrorMsgWrongRequestFormat+": task not pass check", c)
			return
		}
		tid, err := router.processor.TaskStorage.AddTask(task)
		if err != nil {
			taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
			return
		}
		token, _ := util.GenerateTaskAccessToken(tid.Hex(), task.Uid, time.Hour*24)
		requestTaskReply := RequestTaskReply{
			Code: http.StatusOK,
			Msg:  "Create upload task OK",
			Data: TaskResult{
				Type:   "token",
				Result: token,
			},
		}
		c.JSON(http.StatusOK, requestTaskReply)
	case model.DOWNLOAD:
		// req Task 转换为 model Task
		task := RequestTask2Task(&reqTask, model.DOWNLOAD, model.CREATING)
		var taskType model.TaskType
		if reqTask.SourceStoragePlan.StorageMode == "EC" {
			taskType = model.DOWNLOAD
		} else if reqTask.SourceStoragePlan.StorageMode == "Replica" {
			taskType = model.DOWNLOAD_REPLICA
		} else {
			logrus.Warn("wrong storageMode")
			taskRequestReplyErr(util.ErrorCodeWrongStorageType, util.ErrorMsgWrongStorageType, c)
			return
		}
		// 获取文件信息 用于文件流量信息记录
		fileInfo, err := router.processor.FileDatabase.GetFileInfo(task.GetRealSourcePath())
		if err != nil {
			logrus.Errorf("cant get file info: %v%v, err: %v", task.Uid, task.SourcePath, err)
			return
		}

		task.TaskType = taskType
		if task.TaskType == model.DOWNLOAD_REPLICA {
			url, err := router.processor.ProcessGetTmpDownloadUrl(task)
			if err != nil {
				taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
				return
			}
			// 不需要将下载链接写入数据库
			//err = router.processor.WriteDownloadUrlToDB(task, url, task.TaskOptions.SourceStoragePlan.Clouds[0])
			//if err != nil {
			//	logrus.Errorf("write download url to db fail: %v", err)
			//}
			_, err = router.processor.Monitor.AddDownloadTraffic(task.Uid, fileInfo.Size, task.TaskOptions.SourceStoragePlan.Clouds[0])
			if util.Config.EnableHttps { //防止 https 下，浏览器拦截 http 请求
				url = strings.Replace(url, "http://", "https://", 1)
			}
			requestTaskReply := RequestTaskReply{
				Code: http.StatusOK,
				Msg:  "Generate download url OK",
				Data: TaskResult{
					Type:   "url",
					Result: url,
				},
			}
			c.JSON(http.StatusOK, requestTaskReply)
		} else {
			tid, err := router.processor.TaskStorage.AddTask(task)
			if err != nil {
				taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
				return
			}
			requestTaskReply := RequestTaskReply{
				Code: http.StatusOK,
				Msg:  "Generate download task OK",
				Data: TaskResult{
					Type:   "tid",
					Result: tid.Hex(),
				},
			}
			c.JSON(http.StatusOK, requestTaskReply)
		}
	case model.SYNC:
		task := RequestTask2Task(&reqTask, model.SYNC, model.CREATING)
		tid, err := router.processor.TaskStorage.AddTask(task)
		if err != nil {
			taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
			return
		}
		requestTaskReply := RequestTaskReply{
			Code: http.StatusOK,
			Msg:  "Generate sync task OK",
			Data: TaskResult{
				Type:   "tid",
				Result: tid.Hex(),
			},
		}
		c.JSON(http.StatusOK, requestTaskReply)
	case model.DELETE:
		task := RequestTask2Task(&reqTask, model.DELETE, model.CREATING)
		// 删除任务使用同步处理
		err = router.processor.DeleteFileInfo(task)
		if err != nil {
			taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
			return
		}
		tid, err := router.processor.TaskStorage.AddTask(task)
		if err != nil {
			taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
			return
		}
		requestTaskReply := RequestTaskReply{
			Code: http.StatusOK,
			Msg:  "Generate delete task OK",
			Data: TaskResult{
				Type:   "tid",
				Result: tid.Hex(),
			},
		}
		c.JSON(http.StatusOK, requestTaskReply)
	case model.MIGRATE:
		task := RequestTask2Task(&reqTask, model.MIGRATE, model.CREATING)
		tid, err := router.processor.TaskStorage.AddTask(task)
		if err != nil {
			taskRequestReplyErr(util.ErrorCodeInternalErr, err.Error(), c)
			return
		}
		requestTaskReply := RequestTaskReply{
			Code: http.StatusOK,
			Msg:  "Generate delete task OK",
			Data: TaskResult{
				Type:   "tid",
				Result: tid.Hex(),
			},
		}
		c.JSON(http.StatusOK, requestTaskReply)
	default:
		requestTaskReply := RequestTaskReply{
			Code: util.ErrorCodeWrongTaskType,
			Msg:  util.ErrorMsgWrongTaskType,
			Data: TaskResult{},
		}
		c.JSON(util.ErrorCodeWrongTaskType, requestTaskReply)
	}
}

func (router *Router) GetUserTask(c *gin.Context) {
	uid := c.Param("uid")
	tasks := router.processor.TaskStorage.GetUserTask(uid)
	requestTaskReply := RequestGetUserTaskReply{
		Code:  http.StatusOK,
		Msg:   "get user tasks ok",
		Tasks: tasks,
	}
	c.JSON(http.StatusOK, requestTaskReply)
}

func taskRequestReplyErr(errCode int, errMsg string, c *gin.Context) {
	requestTaskReply := RequestTaskReply{
		Code: errCode,
		Msg:  errMsg,
		Data: TaskResult{},
	}
	util.Log(logrus.WarnLevel, "Router task request", errMsg, "", "", "")
	c.JSON(http.StatusBadGateway, requestTaskReply)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, util.Config.LocalCloudID+" JcsPan Transporter version: "+util.GetVersionStr())
}

func (router *Router) FileIndex(c *gin.Context) {
	path := c.Param("path")[1:]
	log.Printf("Index path :%v", path)
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie Fail: %v", err)
	}
	task := model.NewTask(model.INDEX, time.Now(), sidCookie.Value, path, "")
	for obj := range router.processor.ProcessPathIndex(task) {
		fmt.Fprintf(c.Writer, "%s\n", obj.Key)
	}
}

func (router *Router) AddUploadTask(c *gin.Context) {
	//destinationPath := c.Param("path")[1:]
	uid := c.MustGet("tokenUid").(string)
	tid := c.MustGet("tokenTid").(string)
	taskid, err := primitive.ObjectIDFromHex(tid)
	task, err := router.processor.TaskStorage.GetTask(taskid)
	if err != nil {
		log.Printf("Get task fail: %v", err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	destinationPath := task.DestinationPath
	logrus.Infof("upload to :%v", destinationPath)
	// todo: 文件较小时，不落盘，直接内存上传
	c.Request.ParseMultipartForm(32 << 20)
	// 鉴权
	if uid != task.Uid {
		log.Printf("wrong uid")
		http.Error(c.Writer, "wrong uid", http.StatusUnauthorized)
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	randStr := util.GenRandomString(10)
	filePath := util.Config.UploadFileTempPath + handler.Filename + randStr
	// 创建文件，且文件必须不存在
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666) // 此处假设当前目录下已存在test目录
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	sourcePath := filePath

	task.SourcePath = sourcePath
	task.DestinationPath = destinationPath
	task.State = model.WAITING
	// 上传文件后同步写入数据库
	err = router.processor.AddFileInfo(task)
	router.processor.TaskStorage.SetTask(task.Tid, task)
	requestTaskReply := RequestTaskReply{
		Code: http.StatusOK,
		Msg:  "Upload File to Transporter OK",
	}
	c.JSON(http.StatusOK, requestTaskReply)
}

// 获取网盘文件临时下载链接
func (router *Router) GetFile(c *gin.Context) {
	filePath := c.Param("path")[1:]
	log.Printf("get tmp download url: %v", filePath)
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie File: %v", err)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(c.Writer, "Auth Fail")
		return
	}
	task := model.NewTask(model.DOWNLOAD_REPLICA, time.Now(), sidCookie.Value, filePath, "")
	url, err := router.processor.ProcessGetTmpDownloadUrl(task)
	if err != nil {
		log.Printf("Get tmp download url fail: %v", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(c.Writer, "500 ERROR")
	}
	fmt.Fprintln(c.Writer, url)
}

func RequestTask2Task(reqTask *RequestTask, taskType model.TaskType, state model.TaskState) *model.Task {
	var srcCloudsID []string
	var dstCloudsID []string
	for _, cloud := range reqTask.SourceStoragePlan.Clouds {
		srcCloudsID = append(srcCloudsID, cloud.CloudID)
	}
	for _, cloud := range reqTask.DestinationStoragePlan.Clouds {
		dstCloudsID = append(dstCloudsID, cloud.CloudID)
	}
	// 由于传入参数与 transporter 内对于 storagePlan 的 N、K 定义不同，在此需要对 N、K 进行转换
	realSourceN := reqTask.SourceStoragePlan.K
	realSourceK := reqTask.SourceStoragePlan.N - realSourceN
	realDestN := reqTask.DestinationStoragePlan.K
	realDestK := reqTask.DestinationStoragePlan.N - realDestN
	task := model.Task{
		Tid:             primitive.NewObjectID(),
		TaskType:        taskType,
		State:           state,
		StartTime:       time.Time{},
		Uid:             reqTask.Uid,
		SourcePath:      reqTask.SourcePath,
		DestinationPath: reqTask.DestinationPath,
		TaskOptions: &model.TaskOptions{
			SourceStoragePlan: &model.StoragePlan{
				StorageMode: model.StorageModel(reqTask.SourceStoragePlan.StorageMode),
				Clouds:      srcCloudsID,
				N:           realSourceN,
				K:           realSourceK,
			},
			DestinationStoragePlan: &model.StoragePlan{
				StorageMode: model.StorageModel(reqTask.DestinationStoragePlan.StorageMode),
				Clouds:      dstCloudsID,
				N:           realDestN,
				K:           realDestK,
			},
		},
	}
	return &task
}
