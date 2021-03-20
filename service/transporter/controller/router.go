package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Router struct {
	*gin.Engine
	processor TaskProcessor
}

type RequestTask struct {
	TaskType               string             `json:"TaskType"`
	Uid                    string             `json:"Uid"`
	Sid                    string             `json:"Sid"`
	DestinationPath        string             `json:"DestinationPath"`
	SourcePath             string             `json:"SourcePath"`
	SourceStoragePlan      RequestStoragePlan `json:"SourceStoragePlan"`
	DestinationStoragePlan RequestStoragePlan `json:"DestinationStoragePlan"`
}

type RequestStoragePlan struct {
	StorageMode string         `json:"StorageMode"`
	Clouds      []RequestCloud `json:"Clouds"`
	N           int            `json:"N"`
	K           int            `json:"K"`
}

type RequestCloud struct {
	ID string `json:"ID"`
}

func NewRouter(processor TaskProcessor) *Router {
	var router Router
	router = Router{
		Engine:    gin.Default(),
		processor: processor,
	}
	router.GET("/", Index)
	router.POST("/upload/*path", router.AddUploadTask)
	router.GET("/jcspan/*path", router.GetFile)
	router.GET("/index/*path", router.FileIndex)
	router.POST("/task", router.CreateTask)
	router.GET("/cache_file", JWTAuthMiddleware(), router.GetLocalFileByToken)
	rand.Seed(time.Now().Unix())
	return &router
}

func (router *Router) GetLocalFileByToken(c *gin.Context) {
	filePath := c.MustGet("filePath").(string)
	fmt.Println(filePath)
}

func (router *Router) CreateTask(c *gin.Context) {
	var reqTask RequestTask

	err := json.NewDecoder(c.Request.Body).Decode(&reqTask)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch reqTask.TaskType {
	case "Upload":
		var cloudsID []string
		for _, cloud := range reqTask.DestinationStoragePlan.Clouds {
			cloudsID = append(cloudsID, cloud.ID)
		}
		task := model.Task{
			Tid:             primitive.NewObjectID(),
			TaskType:        model.UPLOAD,
			State:           model.BLOCKED,
			StartTime:       time.Time{},
			Uid:             reqTask.Uid,
			SourcePath:      "",
			DestinationPath: reqTask.DestinationPath,
			TaskOptions: &model.TaskOptions{
				SourceStoragePlan: nil,
				DestinationPlan: &model.StoragePlan{
					StorageMode: reqTask.DestinationStoragePlan.StorageMode,
					Clouds:      cloudsID,
					N:           reqTask.DestinationStoragePlan.N,
					K:           reqTask.DestinationStoragePlan.K,
				},
			},
		}
		tid, err := router.processor.taskStorage.AddTask(&task)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		}
		tidStr := tid.Hex()
		fmt.Fprintf(c.Writer, "%v", tidStr)
	case "Download":
		// req Task 转换为 model Task
		var cloudsID []string
		for _, cloud := range reqTask.SourceStoragePlan.Clouds {
			cloudsID = append(cloudsID, cloud.ID)
		}
		var taskType model.TaskType
		if reqTask.SourceStoragePlan.StorageMode == "EC" {
			taskType = model.DOWNLOAD_EC
		} else if reqTask.SourceStoragePlan.StorageMode == "Replica" {
			taskType = model.DOWNLOAD_REPLICA
		} else {
			logrus.Warn("wrong storageMode")
			http.Error(c.Writer, "wrong storage mode", http.StatusBadRequest)

		}

		task := model.Task{
			Tid:             primitive.NewObjectID(),
			TaskType:        taskType,
			State:           model.CREATING,
			StartTime:       time.Time{},
			Uid:             reqTask.Uid,
			SourcePath:      reqTask.SourcePath,
			DestinationPath: "",
			TaskOptions: &model.TaskOptions{
				SourceStoragePlan: &model.StoragePlan{
					StorageMode: reqTask.SourceStoragePlan.StorageMode,
					Clouds:      cloudsID,
					N:           reqTask.SourceStoragePlan.N,
					K:           reqTask.SourceStoragePlan.K,
				},
				DestinationPlan: nil,
			},
		}
		if task.TaskType == model.DOWNLOAD_REPLICA {
			url, err := router.processor.ProcessGetTmpDownloadUrl(&task)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadGateway)
				return
			}
			fmt.Fprintln(c.Writer, url)
		} else {
			tid, err := router.processor.taskStorage.AddTask(&task)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadGateway)
				return
			}
			fmt.Fprintln(c.Writer, tid)
		}

	default:
		http.Error(c.Writer, "wrong task type", http.StatusNotImplemented)
	}
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "JcsPan Transporter")
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
	destinationPath := c.Param("path")[1:]
	log.Printf("upload to :%v", destinationPath)
	sidCookie, err := c.Request.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie Fail: %v", err)
	}
	// todo: 文件较小时，不落盘，直接内存上传
	c.Request.ParseMultipartForm(32 << 20)
	tid := c.Request.FormValue("tid")
	taskid, err := primitive.ObjectIDFromHex(tid)
	task, err := router.processor.taskStorage.GetTask(taskid)
	if err != nil {
		log.Printf("Get task fail: %v", err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	// 鉴权
	user, err := model.Authenticate(sidCookie.Value)
	if err != nil || user.Id != task.Uid {
		log.Printf("wrong sid")
		http.Error(c.Writer, "wrong sid", http.StatusBadGateway)
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// todo: 自定义临时路径
	randStr := genRandomString(10)
	filePath := "./tmp/" + handler.Filename + randStr
	if !isDir("./tmp/") {
		os.MkdirAll("./tmp/", os.ModePerm)
	}
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

	router.processor.taskStorage.SetTask(task.Tid, task)
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

func (router *Router) SimpleSync(c *gin.Context) {
	srcPath := c.Param("srcpath")
	dstPath := c.Param("dstpath")
	sid := c.Request.FormValue("sid")
	router.processor.CreateTask(model.SYNC_SIMPLE, sid, srcPath, dstPath)
	fmt.Println("task simple sync created success")
}

func NewTestRouter(processor TaskProcessor) *Router {
	router := NewRouter(processor)
	router.GET("/test/setcookie", testSetCookie)
	return router
}

func testSetCookie(c *gin.Context) {
	expiration := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:    "sid",
		Value:   "tttteeeesssstttt",
		Expires: expiration,
	}
	http.SetCookie(c.Writer, &cookie)
	fmt.Fprint(c.Writer, "Cookie Already Set")
}
