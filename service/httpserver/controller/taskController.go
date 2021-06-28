package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/service/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserGetTask(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordTaskID:      true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	taskId := (*valueMap)[args.FieldWordTaskID].(string)
	taskIdExist := (*existMap)[args.FieldWordTaskID]
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}

	tasks, getTaskSuccess := dao.TaskDao.GetTask(taskId, userId, taskIdExist)
	if checkDaoSuccess(con, getTaskSuccess) {
		return
	}
	// check if it is correct user
	// TODO
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取任务成功",
		"data": gin.H{
			"Tasks": *tasks,
		},
	})
}
