package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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
	jsi.PUT("/*key", jsi.defaultReply)
	jsi.DELETE("/*key", jsi.defaultReply)
	jsi.GET("/*key", jsi.defaultReply)
	return &jsi
}

func (jsi *JointStorageInterface) GetMethod(c *gin.Context) {
	key := c.Param("key")
	if key == "/" {
		jsi.defaultReply(c)
	} else {
		jsi.defaultReply(c)
	}
}

func (jsi *JointStorageInterface) defaultReply(c *gin.Context) {
	key := c.Param("key")
	c.String(http.StatusNotFound, "method: %s, Path: %s, Key: %s",
		c.Request.Method, c.Request.URL.Path, key)
}

func (jsi *JointStorageInterface) PutObject(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	key := c.Param("key")
	userInfo, err := jsi.processor.UserDatabase.GetUserFromID(uid)
	if err != nil {
		// todo
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

func createTask(uid string, taskType model.TaskType, srcPath string, dstPath string, srcStoragePlan *model.UserStoragePlan,
	dstStoragePlan *model.UserStoragePlan) *model.Task {
	return nil
}
