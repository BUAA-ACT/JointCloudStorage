package newcloud

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/entity"
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
func DaoInit(mongoURI string, databaseMap map[string]*dao.DatabaseConfig) error {
	return dao.NewDao(mongoURI, databaseMap)
}

func IDInit(cid string) {
	localID = cid
}

func PlugIn(mongo string, databasename string, cid string, envMod string, options ...string) error {
	dao.Dao{}
	if err := NewCloudDaoInit(mongo, databasename, cid, envMod); err != nil {
		return err
	}
	for a := range options {

	}
	return nil
}
