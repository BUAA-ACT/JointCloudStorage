package newcloud

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/dao"
)

const (
	CollectionCloud     = "Cloud"
	CollectionTempCloud = "TempCloud"
	CollectionVoteCloud = "VoteCloud"
	CollectionUser      = "User"
	CollectionFile      = "File"
	MigrationAdvice     = "MigrationAdvice"
	codeOK              = 200
	codeBadRequest      = 400
	codeUnauthorized    = 401
	codeInternalError   = 500
)

var (
	localID  string
	errorMsg = map[int]string{
		config.CodeOK:            "OK",
		config.CodeBadRequest:    "Bad Request",
		config.CodeUnauthorized:  "Unauthorized",
		config.CodeInternalError: "Internal Server Error",
	}

	cloudCol     *mongo.Collection
	tempCloudCol *mongo.Collection
	voteCloudCol *mongo.Collection

	env          string
	tempNotFound error
)

func RouteInit(r *gin.Engine) {
	r.POST("/new_cloud", PostNewCloud)
	r.POST("/new_cloud_vote", PostNewCloudVote)
	r.GET("/vote_request", GetVoteRequest)
	r.POST("/cloud_vote", PostCloudVote)
	r.POST("/master_cloud_vote", PostMasterCloudVote)
	r.POST("/cloud_syn", PostCloudSyn)
}

/*
 * NewCloud 的初始化函数，用于初始化mongodb的链接和本地cid
 * mongo：本地mongo数据库地址
 * clouds：database名称
 * cid：本地云的cid
 */
func DaoInit(mongoURI string, databaseMap map[string]map[string]*dao.CollectionConfig) error {
	return dao.NewDao(mongoURI, databaseMap)
}

func IDInit(cid string, envMod string) {
	localID = cid
	env = envMod
}

//func PlugIn(mongo string, databasename string, cid string, envMod string, options ...string) error {
//	dao.Dao{}
//	if err := NewCloudDaoInit(mongo, databasename, cid, envMod); err != nil {
//		return err
//	}
//	for a := range options {
//
//	}
//	return nil
//}
