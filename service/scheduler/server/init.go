package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/dao"
)

const (
	codeOK            = 200
	codeBadRequest    = 400
	codeUnauthorized  = 401
	codeInternalError = 500
	ReplicaMode       = "Replica"
	ECMode            = "EC"
)

var (
	errorMsg = map[int]string{
		codeOK:            "OK",
		codeBadRequest:    "Bad Request",
		codeUnauthorized:  "Unauthorized",
		codeInternalError: "Internal Server Error",
	}
	localId   string
	cloudCol  *mongo.Collection
	userCol   *mongo.Collection
	fileCol   *mongo.Collection
	adviceCol *mongo.Collection
)

func RouteInit(r *gin.Engine) {
	r.GET("/storage_plan", GetStoragePlan)
	r.GET("/download_plan", GetDownloadPlan)
	r.GET("/status", GetStatus)
	r.GET("/all_clouds_status", GetAllCloudsStatus)

	r.POST("/storage_plan", PostStoragePlan)
	r.POST("/metadata", PostMetadata)
	r.POST("/update_clouds", PostUpdateClouds)
}

func DaoInit(mongoURI string, databaseMap map[string]map[string]*dao.CollectionConfig) error {
	return dao.NewDao(mongoURI, databaseMap)
}

func IDInit(cid string) {
	localId = cid
}

func SetCloudCol(thisCol *mongo.Collection) {
	cloudCol = thisCol
}

func SetUserCol(thisCol *mongo.Collection) {
	userCol = thisCol
}

func SetFileCol(thisCol *mongo.Collection) {
	fileCol = thisCol
}

func SetAdviceCol(thisCol *mongo.Collection) {
	adviceCol = thisCol
}
