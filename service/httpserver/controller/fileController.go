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
	userID, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	// check it is file or dir and get files out
	files, listFilesSuccess := dao.FileDao.ListFiles(userID, filePath, tools.IsDir(filePath))
	if !checkDaoSuccess(con, listFilesSuccess) {
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
	//TODO
}

func UserChangeFileName(con *gin.Context) {
	//TODO
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
	userID, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	// database error
	if !checkDaoSuccess(con, infoSuccess) {
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
	preUploadToTransporterResponse, preUploadToTransporterSuccess := transporter.PreUploadFile(filePath, user)
	if !preUploadToTransporterSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析transporter-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// error in transporter
	if preUploadToTransporterResponse.Code != args.CodeOK {
		fmt.Println("transporter fault:")
		fmt.Println("Code: ", preUploadToTransporterResponse.Code)
		fmt.Println("Msg: ", preUploadToTransporterResponse.Msg)
		con.JSON(http.StatusOK, gin.H{
			"code": preUploadToTransporterResponse.Code,
			"msg":  preUploadToTransporterResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// return
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "上传请求已记录",
		"data": gin.H{
			"Token": preUploadToTransporterResponse.Data.Result,
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
	userID, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	//reqId, fileName := regex.FileIdToUserAndFileName(filePath)
	//check file and user
	//if userID != reqId {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDifferentUser,
	//		"msg":  "文件Owner与请求UserId不符合",
	//		"data": gin.H{},
	//	})
	//}
	// check user
	user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	if !checkDaoSuccess(con, infoSuccess) {
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
	// check file status if done -> return url
	files, checkFileSuccess := dao.FileDao.CheckFileStatus(userID, filePath)
	if !checkDaoSuccess(con, checkFileSuccess) {
		return
	}
	var doneURL string
	var doneFlag = false
	for _, file := range *files {
		if file.ReconstructStatus == args.FileReconstructStatusDone {
			doneFlag = true
			doneURL = file.DownloadUrl
		}
	}
	if doneFlag {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeOK,
			"msg":  "重建已完成",
			"data": gin.H{
				"Type":   "url",
				"Result": doneURL,
			},
		})
		return
	}
	// else -> use scheduler's download plan to download file with transporter
	getDownloadPlanResponse, getDownloadPlanSuccess := scheduler.GetDownloadPlanFromScheduler(userID, filePath)
	if !getDownloadPlanSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
	}

	downloadResponse, downloadSuccess := transporter.DownLoadFile(filePath, userID, getDownloadPlanResponse.Data)
	if !downloadSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	// error in transporter
	if downloadResponse.Code != args.CodeOK {
		fmt.Println("transporter fault:")
		fmt.Println("Code: ", downloadResponse.Code)
		fmt.Println("Msg: ", downloadResponse.Msg)
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
	userID, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserHostRole, args.UserGuestRole})
	if !valid {
		return
	}
	user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	// database error
	if !checkDaoSuccess(con, infoSuccess) {
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
	deleteFileFromTransporterSuccess := transporter.DeleteFile(filePath, user)
	if !deleteFileFromTransporterSuccess {
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
