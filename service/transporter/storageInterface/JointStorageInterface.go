package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JointStorageInterface struct {
	*gin.Engine
}

func NewInterface() *JointStorageInterface {
	var jsi JointStorageInterface
	engine := gin.Default()
	engine.Use(util.CORSMiddleware())
	jsi = JointStorageInterface{engine}
	jsi.PUT("/*key", jsi.defaultReply)
	jsi.DELETE("/*key", jsi.defaultReply)
	jsi.GET("/", jsi.defaultReply)
	jsi.GET("/*key", jsi.defaultReply)
	return &jsi
}

func (jsi *JointStorageInterface) defaultReply(c *gin.Context) {
	c.String(http.StatusNotFound, "404 not fund")
}
