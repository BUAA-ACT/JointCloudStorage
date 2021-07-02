package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
	"time"
)

type JointStorageInterface struct {
	*gin.Engine
	processor *controller.TaskProcessor
}

func NewInterface(processor *controller.TaskProcessor) *JointStorageInterface {
	var jsi JointStorageInterface
	engine := gin.Default()
	engine.Use(util.CORSMiddleware())
	jsi = JointStorageInterface{engine, processor}
	jsi.PUT("/*key", JSIAuthMiddleware(), jsi.checkKey, jsi.PutObject)
	jsi.DELETE("/*key", jsi.defaultReply)
	jsi.GET("/*key", JSIAuthMiddleware(), jsi.GetMethod)
	return &jsi
}

func (jsi *JointStorageInterface) GetMethod(c *gin.Context) {
	key := c.Param("key")
	if key == "/" {
		jsi.defaultReply(c)
	} else {
		jsi.checkKey(c)
		jsi.GetObject(c)
	}
}

func (jsi *JointStorageInterface) checkKey(c *gin.Context) {
	key := c.Param("key")
	key = strings.TrimPrefix(key, "/")
	if len(key) == 0 {
		c.String(http.StatusBadRequest, "operate key not found")
		c.Abort()
	}
	c.Set("key", key)
	c.Next()
}

func (jsi *JointStorageInterface) defaultReply(c *gin.Context) {
	key := c.Param("key")
	c.String(http.StatusNotFound, "method: %s, Path: %s, Key: %s",
		c.Request.Method, c.Request.URL.Path, key)
}

func (jsi *JointStorageInterface) PutObject(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	key := c.MustGet("key").(string)
	userInfo, err := jsi.processor.UserDatabase.GetUserFromID(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI PutObject", "get Userinfo fail",
			"", "err", err.Error())
	}
	// 获取用户存储方案
	storagePlan := userInfo.StoragePlan
	fmt.Print(storagePlan)

	// 文件落盘到本地
	f, tempFile := jsi.processor.TempFileStorage.CreateTmpFile(key)
	defer f.Close()
	_, err = io.Copy(f, c.Request.Body)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI PutObject", "copy request body to file fail",
			"", "err", err.Error())
	}

	// 调用 processor ，处理 upload 请求
	task := createTask(userInfo.UserId, model.UPLOAD, tempFile.FilePath, key, nil, &storagePlan)
	err = jsi.processor.ProcessUpload(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI PutObject", "upload task process fail",
			"", "err", err.Error())
	}
	// 返回用户结果
	c.String(http.StatusOK, "")
}

func (jsi *JointStorageInterface) GetObject(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	key := c.MustGet("key").(string)
	userInfo, err := jsi.processor.UserDatabase.GetUserFromID(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI GetObject", "get Userinfo fail",
			"", "err", err.Error())
	}
	// 获取用户存储方案
	storagePlan := userInfo.StoragePlan
	fmt.Print(storagePlan)
	task := createTask(userInfo.UserId, model.DOWNLOAD, key, "", &storagePlan, nil)

	path, err := jsi.processor.RebuildFileToDisk(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI GetObject", "GetObject task process fail",
			"", "err", err.Error())
	}
	c.File(path)
}

func createTask(uid string, taskType model.TaskType, srcPath string, dstPath string, srcStoragePlan *model.UserStoragePlan,
	dstStoragePlan *model.UserStoragePlan) *model.Task {
	var srcCloudsID []string
	var dstCloudsID []string
	var sourceStoragePlan *model.StoragePlan
	if srcStoragePlan != nil {
		for _, cloud := range srcStoragePlan.Clouds {
			srcCloudsID = append(srcCloudsID, cloud.CloudID)
		}
		realSourceN := srcStoragePlan.K
		realSourceK := srcStoragePlan.N - realSourceN
		sourceStoragePlan = &model.StoragePlan{
			StorageMode: model.StorageModel(srcStoragePlan.StorageMode),
			Clouds:      srcCloudsID,
			N:           realSourceN,
			K:           realSourceK,
		}
	}
	var destinationPlan *model.StoragePlan
	if dstStoragePlan != nil {
		for _, cloud := range dstStoragePlan.Clouds {
			dstCloudsID = append(dstCloudsID, cloud.CloudID)
		}
		// 由于传入参数与 transporter 内对于 storagePlan 的 N、K 定义不同，在此需要对 N、K 进行转换
		realDestN := dstStoragePlan.K
		realDestK := dstStoragePlan.N - realDestN
		destinationPlan = &model.StoragePlan{
			StorageMode: model.StorageModel(dstStoragePlan.StorageMode),
			Clouds:      dstCloudsID,
			N:           realDestN,
			K:           realDestK,
		}
	}
	task := model.Task{
		Tid:             primitive.NewObjectID(),
		TaskType:        taskType,
		State:           model.WAITING,
		StartTime:       time.Now(),
		Uid:             uid,
		SourcePath:      srcPath,
		DestinationPath: dstPath,
		TaskOptions: &model.TaskOptions{
			SourceStoragePlan: sourceStoragePlan,
			DestinationPlan:   destinationPlan,
		},
		Progress: 0,
	}
	return &task
}
