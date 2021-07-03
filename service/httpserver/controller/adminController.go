package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AdminAddCloud(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordCloud:       true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	cloud := (*valueMap)[args.FieldWordCloud].(*model.Cloud)
	// check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// check cloud id exist?
	if dao.CloudDao.CheckSameCloudID(cloud.CloudID) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeSameCloudID,
			"msg":  "CloudID已存在,不可重复添加",
			"data": gin.H{},
		})
		return
	}
	// post new cloud to scheduler
	postNewCloudResponse, postNewCloudSuccess := scheduler.PostNewCloudToScheduler(cloud)
	if !postNewCloudSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if postNewCloudResponse.Code != args.CodeOK {
		log.Println("scheduler fault:")
		log.Println("Code: ", postNewCloudResponse.Code)
		log.Println("Msg: ", postNewCloudResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postNewCloudResponse.Code,
			"msg":  postNewCloudResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员新增云节点成功，等待投票",
		"data": gin.H{},
	})
}

func AdminChangeCloudInfo(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordCloud:       true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	cloud := (*valueMap)[args.FieldWordCloud].(*model.Cloud)
	// check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// check cloud id exist?
	if !dao.CloudDao.CheckSameCloudID(cloud.CloudID) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeCloudIDNotExist,
			"msg":  "CloudID不存在",
			"data": gin.H{},
		})
		return
	}
	// post cloud info to scheduler to change and sync it
	postUpdateCloudResponse, postUpdateCloudSuccess := scheduler.PostUpdateCloudToScheduler(cloud)
	if !postUpdateCloudSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if postUpdateCloudResponse.Code != args.CodeOK {
		log.Println("scheduler fault:")
		log.Println("Code: ", postUpdateCloudResponse.Code)
		log.Println("Msg: ", postUpdateCloudResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postUpdateCloudResponse.Code,
			"msg":  postUpdateCloudResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员修改云节点信息成功，等待投票",
		"data": gin.H{},
	})
}

func AdminVoteForCloud(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordCloudID:     true,
		args.FieldWordVoteResult:  true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	cloudID := (*valueMap)[args.FieldWordCloudID].(string)
	voteResult := (*valueMap)[args.FieldWordVoteResult].(bool)
	// check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// post cloud vote result to scheduler
	postCloudVoteResponse, postCloudVoteSuccess := scheduler.PostCloudVoteToScheduler(cloudID, voteResult)
	if !postCloudVoteSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if postCloudVoteResponse.Code != args.CodeOK {
		log.Println("scheduler fault:")
		log.Println("Code: ", postCloudVoteResponse.Code)
		log.Println("Msg: ", postCloudVoteResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postCloudVoteResponse.Code,
			"msg":  postCloudVoteResponse.Msg,
			"data": gin.H{},
		})
		return
	}

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员投票成功",
		"data": gin.H{},
	})
}

func AdminGetVoteRequests(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// get the voting clouds with dao
	voteCloudsResult, voteCloudsSuccess := dao.VoteCloudDao.GetAllVoteCloud()
	if !checkDaoSuccess(con, voteCloudsSuccess) {
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员获取正在投票的云列表成功",
		"data": gin.H{
			"Clouds": *voteCloudsResult,
		},
	})
}

func AdminGetAddedClouds(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// get the added clouds with dao
	addedCloudsResult, addedCloudsSuccess := dao.TempCloudDao.GetAllAddedCloud()
	if !checkDaoSuccess(con, addedCloudsSuccess) {
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员获取已添加的云列表成功",
		"data": gin.H{
			"Clouds": *addedCloudsResult,
		},
	})
}
