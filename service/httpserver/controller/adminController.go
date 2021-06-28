package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminGetAllClouds(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}
	// get clouds with dao
	clouds, getCloudsSuccess := dao.CloudDao.GetAllClouds(userId)
	if !checkDaoSuccess(con, getCloudsSuccess) {
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员获取所有云节点成功",
		"data": gin.H{
			"Clouds": *clouds,
		},
	})
}

func AdminAddCloud(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}

	fmt.Println(userId)

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
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}

	fmt.Println(userId)

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
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAdminRole})
	if !valid {
		return
	}

	fmt.Println(userId)

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "管理员投票成功",
		"data": gin.H{},
	})
}
