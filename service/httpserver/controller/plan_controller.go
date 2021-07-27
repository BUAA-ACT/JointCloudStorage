package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"cloud-storage-httpserver/service/transporter"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserGetAllStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check access token
	user, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}

	// check preference is exist?
	if !user.UserHavePreference() {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePreferenceNotExist,
			"msg":  "请先设置用户偏好",
			"data": gin.H{},
		})
		return
	}
	// get storage plan from scheduler
	response, storagePlanFromSchedulerSuccess := scheduler.GetAllStoragePlanFromScheduler(&user.Preference)
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
	con.Next()
}

func UserGetAdvice(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	user, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	advices, adviceSuccess := dao.MigrationAdviceDao.GetNewAdvice(user.UserID)
	if !checkDaoSuccess(con, adviceSuccess) {
		return
	}
	// choose advice
	empty := make([]model.MigrationAdvice, 0)
	if user.Status != args.UserForbiddenStatus {
		if len(*advices) > 0 && (*advices)[0].Status == args.AdviceStatusChoose {
			advices = &empty
		}
	}
	// return advices
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取建议成功",
		"data": gin.H{
			"Advices": *advices,
		},
	})
	con.Next()
}

func UserAbandonAdvice(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	user, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// delete advice with user status
	if user.Status == args.UserForbiddenStatus {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeStatusForbidden,
			"msg":  "用户正在迁移不能删除",
			"data": gin.H{},
		})
		return
	} else {
		deleteAdviceResult, deleteAdviceSuccess := dao.MigrationAdviceDao.DeleteAdvice(user.UserID)
		if !checkDaoSuccess(con, deleteAdviceSuccess) {
			return
		}
		// nothing has been delete
		if deleteAdviceResult.DeletedCount == 0 {
			con.JSON(http.StatusOK, gin.H{
				"code": args.CodeDeleteNothing,
				"msg":  "重复删除迁移建议",
				"data": gin.H{},
			})
			return
		}
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "抛弃方案成功",
		"data": gin.H{},
	})
	con.Next()
}

func UserChooseStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordStoragePlan: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	storagePlan := (*valueMap)[args.FieldWordStoragePlan].(*model.StoragePlan)
	// check token
	user, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	user, infoSuccess := dao.UserDao.GetUserInfo(user.UserID)
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
	postPlanResponse, postPlanSuccess := scheduler.SetStoragePlanToScheduler(user.UserID, code.AesDecrypt(user.Password, *args.EncryptKey), storagePlan)
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
	credentialSuccess := dao.UserDao.SetUserAccessCredential(user.UserID, &postPlanResponse.Data)
	if !checkDaoSuccess(con, credentialSuccess) {
		return
	}
	// save new plan 不需要再次保存新的存储方案，存储方案交给 scheduler 同步
	//storagePlanSuccess := dao.UserDao.SetUserStoragePlan(userID, storagePlan)
	//if !checkDaoSuccess(con, storagePlanSuccess) {
	//	return
	//}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置存储方案成功",
		"data": gin.H{},
	})
	con.Next()
}

func UserAcceptStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	user, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// check user status
	statusMap := map[string]bool{
		args.UserForbiddenStatus: false,
		args.UserVerifyStatus:    false,
	}
	if !UserCheckStatus(con, user, &statusMap) {
		return
	}
	// forbid user other transportation
	statusSuccess := dao.UserDao.SetUserStatusWithId(user.UserID, args.UserForbiddenStatus)
	if !checkDaoSuccess(con, statusSuccess) {
		return
	}
	// take advice out
	newAdvices, adviceSuccess := dao.MigrationAdviceDao.GetNewAdvice(user.UserID)
	if !checkDaoSuccess(con, adviceSuccess) {
		return
	}
	// advice must be pending
	nowAdvice := (*newAdvices)[0]
	if nowAdvice.Status != args.AdviceStatusPending {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeChangeNothing,
			"msg":  "迁移建议非法,不能进行迁移",
			"data": gin.H{},
		})
	}
	// post to notice scheduler this plan
	postPlanToSchedulerResponse, postPlanToSchedulerSuccess := scheduler.SetStoragePlanToScheduler(user.UserID, code.AesDecrypt(user.Password, *args.EncryptKey), &nowAdvice.StoragePlanNew)
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
		log.Println("scheduler fault:")
		log.Println("Code: ", postPlanToSchedulerResponse.Code)
		log.Println("Msg: ", postPlanToSchedulerResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": postPlanToSchedulerResponse.Code,
			"msg":  postPlanToSchedulerResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// save access credential respond from scheduler
	credentialSuccess := dao.UserDao.SetUserAccessCredential(user.UserID, &postPlanToSchedulerResponse.Data)
	if !checkDaoSuccess(con, credentialSuccess) {
		return
	}
	// save new plan
	storagePlanSuccess := dao.UserDao.SetUserStoragePlan(user.UserID, &nowAdvice.StoragePlanNew)
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
	syncFromTransporterResponse, syncFromTransporterSuccess := transporter.SyncFile("/", user.UserID, sourcePlan, destinationPlan)
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
	// change advice status
	_, setAdviceSuccess := dao.MigrationAdviceDao.SetAdviceStatus(user.UserID, args.AdviceStatusChoose)
	if !checkDaoSuccess(con, setAdviceSuccess) {
		return
	}
	// recover user status : forbidden -> normal ?

	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置存储方案成功",
		"data": gin.H{
			"Transport": true,
			"TaskID":    syncFromTransporterResponse.Data.Result,
		},
	})
	con.Next()
}

/* nonsense

func UserSetStoragePlan(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordStoragePlan: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, fieldRequired)
	if tools.RequiredFieldNotExist(fieldRequired, existMap) {
		return
	}
	accessToken := valueMap[args.FieldWordAccessToken].(string)
	storagePlan := valueMap[args.FieldWordStoragePlan].(*model.StoragePlan)
	// check token
	userID, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	user, success := dao.UserDao.GetUserInfo(userID)
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
	postPlanResponse, postPlanSuccess := scheduler.SetStoragePlanToScheduler(userID, storagePlan)
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
	dao.UserDao.SetUserAccessCredential(userID, &postPlanResponse.Data)
	// save new plan
	dao.UserDao.SetUserStoragePlan(userID, storagePlan)

	// if origin plan exist -> sync it
	if havePlan {
		// use "" to tell transporter migrate all files
		syncResponse, syncSuccess := transporter.SyncFile("", userID, oldPlan, storagePlan)
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
		dao.MigrationAdviceDao.DeleteAdvice(userID)
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
