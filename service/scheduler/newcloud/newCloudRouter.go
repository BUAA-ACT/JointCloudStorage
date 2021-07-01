package newcloud

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, mongo, databasename, cid,envMod string)error{
	if err:=NewCloudInit(mongo,databasename,cid,envMod);err!=nil{
		return err
	}

	r.POST("/new_cloud",PostNewCloud)
	r.POST("/new_cloud_vote",PostNewCloudVote)
	r.GET("/vote_request",GetVoteRequest)
	r.POST("/cloud_vote",PostCloudVote)
	r.POST("/master_cloud_vote",PostMasterCloudVote)
	r.POST("/cloud_syn",PostCloudSyn)

	return nil
}