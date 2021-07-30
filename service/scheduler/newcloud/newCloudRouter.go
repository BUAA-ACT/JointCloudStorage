package newcloud

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/entity"
)

func router(r *gin.Engine) {
	r.POST("/new_cloud", PostNewCloud)
	r.POST("/new_cloud_vote", PostNewCloudVote)
	r.GET("/vote_request", GetVoteRequest)
	r.POST("/cloud_vote", PostCloudVote)
	r.POST("/master_cloud_vote", PostMasterCloudVote)
	r.POST("/cloud_syn", PostCloudSyn)
}

func PlugIn(mongo, databasename, cid, envMod string) error {
	dao.Dao{}
	if err := NewCloudInit(mongo, databasename, cid, envMod); err != nil {
		return err
	}

	return nil
}
