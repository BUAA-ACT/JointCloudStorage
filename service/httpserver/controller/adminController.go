package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllClouds(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	_, userRole, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// get clouds from scheduler
	getCloudsResponse, getCloudsSuccess := scheduler.GetAllCloudsFromScheduler()
	if !getCloudsSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if getCloudsResponse.Code != args.CodeOK {
		con.JSON(http.StatusOK, gin.H{
			"code": getCloudsResponse.Code,
			"msg":  getCloudsResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// if user is not admin then cover clouds ak&sk
	if userRole != args.UserAdminRole {
		for _, cloud := range getCloudsResponse.Data {
			cloud.SecretKey = ""
			cloud.AccessKey = ""
		}
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取所有云节点成功",
		"data": gin.H{
			"Clouds": getCloudsResponse.Data,
		},
	})
}

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
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", postNewCloudResponse.Code)
		fmt.Println("Msg: ", postNewCloudResponse.Msg)
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
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", postUpdateCloudResponse.Code)
		fmt.Println("Msg: ", postUpdateCloudResponse.Msg)
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
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", postCloudVoteResponse.Code)
		fmt.Println("Msg: ", postCloudVoteResponse.Msg)
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
	// get the voting clouds from scheduler
	getVoteRequestsResponse, getVoteRequestsSuccess := scheduler.GetVoteRequestsFromScheduler()
	if !getVoteRequestsSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if getVoteRequestsResponse.Code != args.CodeOK {
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", getVoteRequestsResponse.Code)
		fmt.Println("Msg: ", getVoteRequestsResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": getVoteRequestsResponse.Code,
			"msg":  getVoteRequestsResponse.Msg,
			"data": gin.H{},
		})
		return
	}

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员获取正在投票的云列表成功",
		"data": gin.H{
			"Clouds": getVoteRequestsResponse.Data,
		},
	})
}
