package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"cloud-storage-httpserver/service/transporter"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserGetFiles(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordFilePath:    true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	filePath := (*valueMap)[args.FieldWordFilePath].(string)
	//check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	// check it is file or dir and get files out
	files, success := dao.FileDao.ListFiles(userId, filePath, tools.IsDir(filePath))
	if !success {
		// database err
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	//if len(*files) == 0 {
	//	// don't have such files
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeFieldNotExist,
	//		"msg":  "未找到文件",
	//		"data": gin.H{},
	//	})
	//	return
	//}
	var integrateFiles []model.FileAndDir = make([]model.FileAndDir, 0)
	dirMap := make(map[string]model.File)
	for _, file := range *files {
		remainString := strings.TrimPrefix(file.FileName, filePath)
		segments := strings.Split(remainString, "/")
		if len(segments) <= 1 {
			// is file -> add it
			newReturnFile := model.FileAndDir{
				FileType: args.FileTypeFile,
				FileInfo: file,
			}
			integrateFiles = append(integrateFiles, newReturnFile)
		} else {
			// is dir
			dirName := segments[0]
			wholeDirName := filePath + dirName + "/"
			dirFile, ok := dirMap[dirName]
			if ok {
				// has recorded dir -> add size
				dirFile.Size += file.Size
			} else {
				// record dir
				newDirFile := model.File{
					FileID:   file.Owner + wholeDirName,
					Owner:    file.Owner,
					Size:     file.Size,
					FileName: wholeDirName,
				}
				dirMap[dirName] = newDirFile
			}
		}
	}
	for _, file := range dirMap {
		// add dir
		newReturnFile := model.FileAndDir{
			FileType: args.FileTypeFile,
			FileInfo: file,
		}
		integrateFiles = append(integrateFiles, newReturnFile)
	}
	// return ok
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取文件成功",
		"data": gin.H{
			"Files": integrateFiles,
		},
	})

}

func UserChangeFilePath(con *gin.Context) {

}

func UserChangeFileName(con *gin.Context) {

}

func UserPreUploadFile(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordFilePath:    true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	filePath := (*valueMap)[args.FieldWordFilePath].(string)
	//check access token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	user, success := dao.UserDao.GetUserInfo(userId)
	// database error
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	// check have storage plan?
	if !user.UserHaveStoragePlan() {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeStoragePlanNotExist,
			"msg":  "用户存储方案不存在",
			"data": gin.H{},
		})
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

	// preUpload and get a token
	response, success := transporter.PreUploadFile(filePath, user)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析transporter-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// error in transporter
	if response.Code != args.CodeOK {
		con.JSON(http.StatusOK, gin.H{
			"code": response.Code,
			"msg":  response.Msg,
			"data": gin.H{},
		})
		return
	}
	// return
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "上传请求已记录",
		"data": gin.H{
			"Token": response.Data.Result,
		},
	})
}

func UserDownloadFile(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordFilePath:    true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	filePath := (*valueMap)[args.FieldWordFilePath].(string)
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	//reqId, fileName := regex.FileIdToUserAndFileName(filePath)
	//check file and user
	//if userId != reqId {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDifferentUser,
	//		"msg":  "文件Owner与请求UserId不符合",
	//		"data": gin.H{},
	//	})
	//}
	// check user
	user, success := dao.UserDao.GetUserInfo(userId)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	if !user.UserHaveStoragePlan() {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeStoragePlanNotExist,
			"msg":  "用户存储方案不存在",
			"data": gin.H{},
		})
	}
	// check user status
	statusMap := map[string]bool{
		args.UserForbiddenStatus: false,
		args.UserVerifyStatus:    false,
	}
	if !UserCheckStatus(con, user, &statusMap) {
		return
	}

	// use scheduler's download plan to download file with transporter
	getDownloadPlanResponse, success := scheduler.GetDownloadPlanFromScheduler(userId, filePath)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
	}

	downloadResponse, success := transporter.DownLoadFile(filePath, userId, getDownloadPlanResponse.Data)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// error in transporter
	if downloadResponse.Code != args.CodeOK {
		con.JSON(http.StatusOK, gin.H{
			"code": downloadResponse.Code,
			"msg":  downloadResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// return data
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "下载请求已送达",
		"data": downloadResponse.Data,
	})
}

func UserDeleteFile(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordFilePath:    true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	filePath := (*valueMap)[args.FieldWordFilePath].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	user, success := dao.UserDao.GetUserInfo(userId)
	// database error
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	// check have storage plan?
	if !user.UserHaveStoragePlan() {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeStoragePlanNotExist,
			"msg":  "用户存储方案不存在",
			"data": gin.H{},
		})
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

	// delete file with transporter
	success = transporter.DeleteFile(filePath, user)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析transporter-json信息有误",
			"data": gin.H{},
		})
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "删除文件成功",
		"data": gin.H{},
	})

}
