package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
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
