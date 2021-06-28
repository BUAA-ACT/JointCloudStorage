package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"cloud-storage-httpserver/service/transporter"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserGetAllStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check access token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	thisUser, infoSuccess := dao.UserDao.GetUserInfo(userId)
	if !checkDaoSuccess(con, infoSuccess) {
		return
	}
	// check preference is exist?
	if !thisUser.UserHavePreference() {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePreferenceNotExist,
			"msg":  "请先设置用户偏好",
			"data": gin.H{},
		})
		return
	}
	// get storage plan from scheduler
	response, storagePlanFromSchedulerSuccess := scheduler.GetAllStoragePlanFromScheduler(thisUser.Preference)
	if !storagePlanFromSchedulerSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// wrong in scheduler
	if response.Code != args.CodeOK {
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", response.Code)
		fmt.Println("Msg: ", response.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": response.Code,
			"msg":  response.Msg,
			"data": gin.H{},
		})
		return
	}
	// return ok
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取存储方案成功",
		"data": gin.H{
			"RequestID":         response.RequestID,
			"TrafficPriceFirst": response.Data.TrafficPriceFirst,
			"StoragePriceFirst": response.Data.StoragePriceFirst,
		},
	})
}

func UserGetAdvice(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	advices, adviceSuccess := dao.MigrationAdviceDao.GetNewAdvice(userId)
	if !checkDaoSuccess(con, adviceSuccess) {
		return
	}
	fmt.Print("advices: ")
	fmt.Println(*advices)
	// return advices
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取建议成功",
		"data": gin.H{
			"Advices": *advices,
		},
	})
	return
}

func UserAbandonAdvice(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	deleteAdviceResult, deleteAdviceSuccess := dao.MigrationAdviceDao.DeleteAdvice(userId)
	if !checkDaoSuccess(con, deleteAdviceSuccess) {
		return
	}
	// nothing has been delete
	if deleteAdviceResult.DeletedCount == 0 {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDeleteNothing,
			"msg":  "没有删除任何方案",
			"data": gin.H{},
		})
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "抛弃方案成功",
		"data": gin.H{},
	})
	return
}

func UserChooseStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordStoragePlan: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	storagePlan := (*valueMap)[args.FieldWordStoragePlan].(*model.StoragePlan)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	user, infoSuccess := dao.UserDao.GetUserInfo(userId)
	if !checkDaoSuccess(con, infoSuccess) {
		return
	}
	// check there is origin plan
	oldPlan := &user.StoragePlan
	if oldPlan.N > 0 {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeAlreadyHaveStoragePlan,
			"msg":  "重复初始化存储方案",
			"data": gin.H{},
		})
		return
	}
	// post to notice scheduler this plan
	postPlanResponse, postPlanSuccess := scheduler.SetStoragePlanToScheduler(userId, storagePlan)
	if !postPlanSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postPlanResponse.Code != args.CodeOK {
		// error in scheduler
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", postPlanResponse.Code)
		fmt.Println("Msg: ", postPlanResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postPlanResponse.Code,
			"msg":  postPlanResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// save access credential respond from scheduler
	credentialSuccess := dao.UserDao.SetUserAccessCredential(userId, &postPlanResponse.Data)
	if !checkDaoSuccess(con, credentialSuccess) {
		return
	}
	// save new plan
	storagePlanSuccess := dao.UserDao.SetUserStoragePlan(userId, storagePlan)
	if !checkDaoSuccess(con, storagePlanSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置存储方案成功",
		"data": gin.H{},
	})
}

func UserAcceptStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	// check user status
	user, infoSuccess := dao.UserDao.GetUserInfo(userId)
	if !checkDaoSuccess(con, infoSuccess) {
		return
	}
	statusMap := map[string]bool{
		args.UserForbiddenStatus: false,
		args.UserVerifyStatus:    false,
	}
	if !UserCheckStatus(con, user, &statusMap) {
		return
	}
	// forbid user other transportation
	statusSuccess := dao.UserDao.SetUserStatusWithId(userId, args.UserForbiddenStatus)
	if !checkDaoSuccess(con, statusSuccess) {
		return
	}
	// take advice out
	newAdvices, adviceSuccess := dao.MigrationAdviceDao.GetNewAdvice(userId)
	if !checkDaoSuccess(con, adviceSuccess) {
		return
	}
	nowAdvice := (*newAdvices)[0]
	// post to notice scheduler this plan
	postPlanToSchedulerResponse, postPlanToSchedulerSuccess := scheduler.SetStoragePlanToScheduler(userId, &nowAdvice.StoragePlanNew)
	// save new plan
	if !postPlanToSchedulerSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postPlanToSchedulerResponse.Code != args.CodeOK {
		// error in scheduler
		fmt.Println("scheduler fault:")
		fmt.Println("Code: ", postPlanToSchedulerResponse.Code)
		fmt.Println("Msg: ", postPlanToSchedulerResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postPlanToSchedulerResponse.Code,
			"msg":  postPlanToSchedulerResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// save access credential respond from scheduler
	credentialSuccess := dao.UserDao.SetUserAccessCredential(userId, &postPlanToSchedulerResponse.Data)
	if !checkDaoSuccess(con, credentialSuccess) {
		return
	}
	// save new plan
	storagePlanSuccess := dao.UserDao.SetUserStoragePlan(userId, &nowAdvice.StoragePlanNew)
	if !checkDaoSuccess(con, storagePlanSuccess) {
		return
	}
	// prepare the plan that transporter need
	sourcePlan := &model.StoragePlan{
		StorageMode: "Migrate",
		Clouds:      nowAdvice.CloudsOld,
	}
	destinationPlan := &model.StoragePlan{
		StorageMode: "Migrate",
		Clouds:      nowAdvice.CloudsNew,
	}

	// use "" to tell transporter migrate all files
	syncFromTransporterResponse, syncFromTransporterSuccess := transporter.SyncFile("", userId, sourcePlan, destinationPlan)
	if !syncFromTransporterSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析transporter-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// error in transporter
	if syncFromTransporterResponse.Code != args.CodeOK {
		con.JSON(http.StatusOK, gin.H{
			"code": syncFromTransporterResponse.Code,
			"msg":  syncFromTransporterResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// delete advice
	deleteAdviceResult, deleteAdviceSuccess := dao.MigrationAdviceDao.DeleteAdvice(userId)
	if !checkDaoSuccess(con, deleteAdviceSuccess) {
		return
	}
	// recover user status : forbidden -> normal ?

	// delete nothing
	if deleteAdviceResult.DeletedCount == 0 {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDeleteNothing,
			"msg":  "没有删除任何方案",
			"data": gin.H{},
		})
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置存储方案成功",
		"data": gin.H{
			"Transport": true,
			"TaskID":    syncFromTransporterResponse.Data.Result,
		},
	})

}

/* nonsense

func UserSetStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordStoragePlan: true,
	}
	valueMap, existMap := getQueryAndReturn(con, fieldRequired)
	if tools.RequiredFieldNotExist(fieldRequired, existMap) {
		return
	}
	accessToken := valueMap[args.FieldWordAccessToken].(string)
	storagePlan := valueMap[args.FieldWordStoragePlan].(*model.StoragePlan)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	user, success := dao.UserDao.GetUserInfo(userId)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	// check there is origin plan
	oldPlan := &user.StoragePlan
	havePlan := oldPlan.N > 0
	// post to notice scheduler this plan
	postPlanResponse, postPlanSuccess := scheduler.SetStoragePlanToScheduler(userId, storagePlan)
	// save new plan
	if !postPlanSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postPlanResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": postPlanResponse.Code,
			"msg":  postPlanResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// save access credential respond from scheduler
	dao.UserDao.SetUserAccessCredential(userId, &postPlanResponse.Data)
	// save new plan
	dao.UserDao.SetUserStoragePlan(userId, storagePlan)

	// if origin plan exist -> sync it
	if havePlan {
		// use "" to tell transporter migrate all files
		syncResponse, syncSuccess := transporter.SyncFile("", userId, oldPlan, storagePlan)
		if !syncSuccess {
			con.JSON(http.StatusOK, gin.H{
				"code": args.CodeJsonError,
				"msg":  "解析transporter-json信息有误",
				"data": gin.H{},
			})
			return
		}
		// error in transporter
		if syncResponse.Code != args.CodeOK {
			con.JSON(http.StatusOK, gin.H{
				"code": syncResponse.Code,
				"msg":  syncResponse.Msg,
				"data": gin.H{},
			})
			return
		}
		//delete advice
		dao.MigrationAdviceDao.DeleteAdvice(userId)
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeOK,
			"msg":  "设置存储方案成功",
			"data": gin.H{
				"Transport": true,
				"TaskID":    syncResponse.Data.Result,
			},
		})
	} else {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeOK,
			"msg":  "设置存储方案成功",
			"data": gin.H{
				"Transport": false,
			},
		})
	}
}

*/
