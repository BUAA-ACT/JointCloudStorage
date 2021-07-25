package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"errors"
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
	object := jsi.Group("/object")
	{
		object.PUT("/*key", jsi.JSIAuthMiddleware(), jsi.checkKey, jsi.PutObject)
		object.DELETE("/*key", jsi.JSIAuthMiddleware(), jsi.checkKey, jsi.DeleteObject)
		object.GET("/*key", jsi.JSIAuthMiddleware(), jsi.GetMethod)
	}
	state := jsi.Group("/state")
	{
		state.GET("/storage", jsi.JSIAuthMiddleware(), jsi.GetStorageInfo)
		state.GET("/plan", jsi.JSIAuthMiddleware(), jsi.GetStoragePlan)
		state.POST("/plan", jsi.JSIAuthMiddleware(), jsi.PostStoragePlan)
		state.GET("/server", jsi.GetServerInfo)
	}
	jsi.GET("/", jsi.GetServerInfo)

	return &jsi
}

func (jsi *JointStorageInterface) GetMethod(c *gin.Context) {
	key := c.Param("key")
	if key == "/" {
		jsi.GetObjectList(c)
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

type ServerInfo struct {
	CloudId            string
	TransporterVersion string
	ServerTime         time.Time
}

func (jsi *JointStorageInterface) GetServerInfo(c *gin.Context) {
	info := ServerInfo{
		CloudId:            util.Config.LocalCloudID,
		TransporterVersion: util.GetVersionStr(),
		ServerTime:         time.Now(),
	}
	c.JSON(http.StatusOK, &info)
}

func (jsi *JointStorageInterface) PutObject(c *gin.Context) {
	key := c.MustGet("key").(string)
	userInfo := c.MustGet("userInfo").(*model.User)
	// 获取用户存储方案
	storagePlan := userInfo.StoragePlan
	fmt.Print(storagePlan)

	// 文件落盘到本地
	f, tempFile := jsi.processor.TempFileStorage.CreateTmpFile(key)
	defer f.Close()
	_, err := io.Copy(f, c.Request.Body)
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
		c.String(http.StatusInternalServerError, err.Error())
	}
	// 返回用户结果
	c.String(http.StatusOK, "")
}

func (jsi *JointStorageInterface) GetObject(c *gin.Context) {
	key := c.MustGet("key").(string)
	userInfo := c.MustGet("userInfo").(*model.User)
	// 获取用户存储方案
	storagePlan := userInfo.StoragePlan
	fmt.Print(storagePlan)
	task := createTask(userInfo.UserId, model.DOWNLOAD, key, "", &storagePlan, nil)

	path, err := jsi.processor.RebuildFileToDisk(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI GetObject", "GetObject task process fail",
			"", "err", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.File(path)
}

func (jsi *JointStorageInterface) DeleteObject(c *gin.Context) {
	key := c.MustGet("key").(string)
	userInfo := c.MustGet("userInfo").(*model.User)
	// 获取用户存储方案
	storagePlan := userInfo.StoragePlan
	fmt.Print(storagePlan)
	task := createTask(userInfo.UserId, model.DELETE, key, "", &storagePlan, nil)
	var err error
	err = jsi.processor.DeleteFileInfo(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI DeleteObject", "processor delete file info fail",
			"", "err", err.Error())
	}
	err = jsi.processor.DeleteStorageFile(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI DeleteObject", "processor delete fail",
			"", "err", err.Error())
	}
}

func (jsi *JointStorageInterface) GetObjectList(c *gin.Context) {
	prefix := c.Query("keyPrefix")
	userInfo := c.MustGet("userInfo").(*model.User)
	task := createTask(userInfo.UserId, model.INDEX, prefix, "", nil, nil)
	files, err := jsi.processor.ProcessIndexFile(task)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI GetObjectList", "get Userinfo fail",
			"", "err", err.Error())
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusOK, files)
}

func (jsi *JointStorageInterface) GetStorageInfo(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(*model.User)
	c.JSON(http.StatusOK, userInfo.DataStats)
}

func (jsi *JointStorageInterface) GetStoragePlan(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(*model.User)
	c.JSON(http.StatusOK, userInfo.StoragePlan)
}

func (jsi *JointStorageInterface) PostStoragePlan(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(*model.User)
	var plan model.StoragePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid storage plan"))
		return
	}
	userPlan, err := jsi.convertToUserStoragePlan(plan)
	if err != nil {
		util.Log(logrus.ErrorLevel, "JSI Post StoragePlan", "convert StoragePlan fail",
			"", "err", err.Error())
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid storage plan"))
		return
	}
	err = jsi.processor.Scheduler.SetUserStoragePlan(userInfo, userPlan)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("set storage plan fail"))
		return
	}
	task := createTask(userInfo.UserId, model.SYNC, "/", "/", &userInfo.StoragePlan, userPlan)
	tid, err := jsi.processor.AddTask(task)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("internal error"))
		return
	}
	c.String(http.StatusOK, tid.String())
}

func (jsi *JointStorageInterface) convertToUserStoragePlan(plan model.StoragePlan) (*model.UserStoragePlan, error) {
	var clouds []model.Cloud
	for _, c := range plan.Clouds {
		cloud, err := jsi.processor.CloudDatabase.GetCloudInfoFromCloudID(c)
		if err != nil {
			return nil, errors.New("get cloud info from cloud id fail")
		}
		clouds = append(clouds, *cloud)
	}
	return &model.UserStoragePlan{
		N:            plan.N,
		K:            plan.K,
		StorageMode:  string(plan.StorageMode),
		Clouds:       clouds,
		StoragePrice: 0,
		TrafficPrice: 0,
		Availability: 0,
	}, nil
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
			SourceStoragePlan:      sourceStoragePlan,
			DestinationStoragePlan: destinationPlan,
		},
		Progress: 0,
	}
	return &task
}
