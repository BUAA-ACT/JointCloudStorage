package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/service/tools"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func UserGetTask(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordTaskID:      true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	taskId := (*valueMap)[args.FieldWordTaskID].(string)
	taskIdExist := (*existMap)[args.FieldWordTaskID]
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}

	tasks, getTaskSuccess := dao.TaskDao.GetTask(taskId, userID, taskIdExist)
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

func UserGetMigrationTask(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	ws, weErr := upGrader.Upgrade(con.Writer, con.Request, nil)
	if tools.PrintError(weErr) {
		return
	}
	valueMap, existMap := getVerifyAndReturnWithWebSocket(ws, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// get the task
	go func() {
		var returnMap gin.H
		var progressNow float64 = 0
		for {
			// get task with dao
			task, ok := dao.TaskDao.GetUserMigrate(userID)
			if !ok {
				returnMap = gin.H{
					"code": args.CodeDatabaseError,
					"msg":  "数据库错误",
					"data": gin.H{},
				}
				writeErr := ws.WriteJSON(returnMap)
				if writeErr != nil {
					log.Println("fucking writing json problem: " + writeErr.Error())
				}
				break
			}
			// if it is done
			if task.TaskState == args.TaskStateFailed || task.TaskState == args.TaskStateFinished || task.TaskState == args.TaskStateBlocked {
				returnMap = gin.H{
					"code": args.CodeOK,
					"msg":  "已完成",
					"data": task,
				}
				writeErr := ws.WriteJSON(returnMap)
				if writeErr != nil {
					log.Println("fucking writing json problem: " + writeErr.Error())
				}
				break
			}
			if task.Progress >= progressNow && task.TaskState != args.TaskStateCreating && task.TaskState != args.TaskStateWaiting {
				// return the progress of the task
				returnMap = gin.H{
					"code": args.CodeOK,
					"msg":  "正在同步中",
					"data": task,
				}
				writeErr := ws.WriteJSON(returnMap)
				if writeErr != nil {
					log.Println("fucking writing json problem: " + writeErr.Error())
				}
			}
			// sleep in 0.1s
			time.Sleep(100000000)
		}
	}()
}
