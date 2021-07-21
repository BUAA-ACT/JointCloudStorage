package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
	"cloud-storage-httpserver/service/scheduler"
	"cloud-storage-httpserver/service/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func UserRegister(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordEmail:    true,
		args.FieldWordPassword: true,
		args.FieldWordNickname: false,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	email := (*valueMap)[args.FieldWordEmail].(string)
	password := (*valueMap)[args.FieldWordPassword].(string)
	var nickname string
	if (*existMap)[args.FieldWordNickname] {
		nickname = (*valueMap)[args.FieldWordNickname].(string)
	} else {
		nickname = ""
	}

	// check same email
	if dao.UserDao.CheckSameEmail(email) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeSameEmail,
			"msg":  "用户邮箱已被注册",
			"data": gin.H{},
		})
		return
	}
	nowTime := time.Now()
	// uuid or email?
	userID := email
	// save with dao
	user := &model.User{
		UserID:       userID,
		Email:        email,
		Password:     code.AesEncrypt(password, *args.EncryptKey),
		Nickname:     nickname,
		Role:         args.UserHostRole,
		Avatar:       "default-avatar.png",
		CreateTime:   nowTime,
		LastModified: nowTime,
		Status:       args.UserNormalStatus,
	}
	// save user in db
	userSuccess := dao.UserDao.CreateNewUser(*user)
	if !checkDaoSuccess(con, userSuccess) {
		return
	}
	// record verify code
	verifyCode := code.GenVerifyCode()
	verifyCodeSuccess := dao.VerifyCodeDao.InsertVerifyCode(email, verifyCode)
	if !checkDaoSuccess(con, verifyCodeSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户注册已记录",
		"data": gin.H{},
	})
}

func UserCheckVerifyCode(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordEmail:      true,
		args.FieldWordVerifyCode: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	email := (*valueMap)[args.FieldWordEmail].(string)
	verifyCode := (*valueMap)[args.FieldWordVerifyCode].(string)
	// check email
	if !dao.UserDao.CheckSameEmail(email) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeEmailNotExist,
			"msg":  "用户还未注册",
			"data": gin.H{},
		})
		return
	}
	// check code

	verifyEmailSuccess := dao.VerifyCodeDao.VerifyEmail(email, verifyCode)
	if verifyEmailSuccess {
		// update user status
		changeStatusSuccess := dao.UserDao.SetUserStatusWithEmail(email, args.UserNormalStatus)
		if !checkDaoSuccess(con, changeStatusSuccess) {
			return
		}
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeOK,
			"msg":  "验证码正确",
			"data": gin.H{},
		})

	} else {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeVerifyFail,
			"msg":  "验证码错误",
			"data": gin.H{},
		})
	}
}

func UserLogin(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordEmail:    true,
		args.FieldWordPassword: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	email := (*valueMap)[args.FieldWordEmail].(string)
	password := (*valueMap)[args.FieldWordPassword].(string)
	// check email exist
	if !dao.UserDao.CheckSameEmail(email) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeEmailNotExist,
			"msg":  "用户邮箱不存在",
			"data": gin.H{},
		})
		return
	}
	// check password
	user, loginSuccess := dao.UserDao.LoginWithEmail(email, password)
	if !loginSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "密码错误",
			"data": gin.H{},
		})
		return
	}
	// check user status
	statusMap := map[string]bool{
		args.UserVerifyStatus: false,
	}
	if !UserCheckStatus(con, user, &statusMap) {
		return
	}
	// gen token
	token := code.GenToken().String()
	tokenSuccess := dao.AccessTokenDao.InsertAccessToken(token, user.UserID)
	if !checkDaoSuccess(con, tokenSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "登录成功",
		"data": gin.H{
			"AccessToken": token,
		},
	})
}

func UserLogout(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token is valid
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// delete token
	deleteTokenResult, deleteTokenSuccess := dao.AccessTokenDao.DeleteAccessToken(accessToken)
	if !checkDaoSuccess(con, deleteTokenSuccess) {
		return
	}
	// no token been removed
	if deleteTokenResult.DeletedCount == 0 {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDeleteNothing,
			"msg":  "¿已经退出过了啊¿",
			"data": gin.H{},
		})
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "白白了您呐!",
		"data": gin.H{},
	})
}

func UserCheckValidity(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	//check token
	_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// valid
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "令牌有效",
		"data": gin.H{
			"AccessToken": accessToken,
		},
	})
}

func UserChangePassword(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken:    true,
		args.FieldWordOriginPassword: true,
		args.FieldWordNewPassword:    true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	originPassword := (*valueMap)[args.FieldWordOriginPassword].(string)
	newPassword := (*valueMap)[args.FieldWordNewPassword].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// verify role
	user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	if !checkDaoSuccess(con, infoSuccess) {
		return
	}
	// it can be modify into UserRoles ↓
	if user.Role == args.UserAllRole {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "Guest禁止修改密码",
			"data": gin.H{},
		})
		return
	}

	// verify origin password
	loginSuccess := dao.UserDao.LoginWithId(userID, originPassword)
	if !loginSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "原密码错误",
			"data": gin.H{},
		})
		return
	}
	// dao change password
	changePasswordSuccess := dao.UserDao.SetUserPassword(userID, newPassword)
	if !checkDaoSuccess(con, changePasswordSuccess) {
		return
	}
	// let scheduler change and sync password
	storagePlan := &user.StoragePlan
	postPlanResponse, postPlanSuccess := scheduler.SetStoragePlanToScheduler(userID, newPassword, storagePlan)
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
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "修改密码成功",
		"data": gin.H{},
	})
}

func UserChangeEmail(con *gin.Context) {
	// we can't change email now
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordNewEmail:    true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	newEmail := (*valueMap)[args.FieldWordNewEmail].(string)
	// check token
	userID, userRole, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// can't change email now!
	//user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	//if !checkDaoSuccess(con, infoSuccess) {
	//	return
	//}
	// it can be modify into roles
	if userRole == args.UserHostRole || userRole == args.UserGuestRole {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "禁止修改邮箱",
			"data": gin.H{},
		})
		return
	}

	// check same email
	if dao.UserDao.CheckSameEmail(newEmail) {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeSameEmail,
			"msg":  "用户邮箱已被注册",
			"data": gin.H{},
		})
		return
	}
	// save with dao
	emailSuccess := dao.UserDao.SetUserEmail(userID, newEmail)
	if !checkDaoSuccess(con, emailSuccess) {
		return
	}
	statusSuccess := dao.UserDao.SetUserStatusWithId(userID, args.UserVerifyStatus)
	if !checkDaoSuccess(con, statusSuccess) {
		return
	}
	// verify code
	verifyCode := code.GenVerifyCode()
	insertVerifyCodeSuccess := dao.VerifyCodeDao.InsertVerifyCode(newEmail, verifyCode)
	if !checkDaoSuccess(con, insertVerifyCodeSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户更改邮箱已记录",
		"data": gin.H{},
	})
}

func UserChangeNickname(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordNewNickname: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	newNickname := (*valueMap)[args.FieldWordNickname].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	nicknameSuccess := dao.UserDao.SetUserNickname(userID, newNickname)
	if !checkDaoSuccess(con, nicknameSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "修改昵称成功",
		"data": gin.H{},
	})
}

func UserGetInfo(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	user, infoSuccess := dao.UserDao.GetUserInfo(userID)
	if !checkDaoSuccess(con, infoSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "获取信息成功",
		"data": gin.H{
			"UserInfo": user,
		},
	})
}

func UserSetPreference(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken:  true,
		args.FieldWordVendor:       true,
		args.FieldWordStoragePrice: true,
		args.FieldWordTrafficPrice: true,
		args.FieldWordAvailability: true,
		args.FieldWordLatency:      false,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	vendor := (*valueMap)[args.FieldWordVendor].(uint64)
	storagePrice := (*valueMap)[args.FieldWordStoragePrice].(float64)
	trafficPrice := (*valueMap)[args.FieldWordTrafficPrice].(float64)
	availability := (*valueMap)[args.FieldWordAvailability].(float64)
	var latency *map[string]uint64
	if (*existMap)[args.FieldWordLatency] {
		latency = (*valueMap)[args.FieldWordLatency].(*map[string]uint64)
	} else {
		latency = &map[string]uint64{}
	}
	//check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	preference := &model.Preference{
		Vendor:       vendor,
		StoragePrice: storagePrice,
		TrafficPrice: trafficPrice,
		Availability: availability,
		Latency:      *latency,
	}
	// set preference
	preferenceSuccess := dao.UserDao.SetUserPreference(userID, preference)
	if !checkDaoSuccess(con, preferenceSuccess) {
		return
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置个人偏好成功",
		"data": gin.H{},
	})
}

func UserGetKeys(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// get all keys with userID
	keys, getKeysSuccess := dao.AccessKeyDao.GetAllKeys(userID)
	if !checkDaoSuccess(con, getKeysSuccess) {
		return
	}
	// success with gen keys
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户获取所有密钥成功",
		"data": gin.H{
			"Keys": *keys,
		},
	})
}

func UserAddKey(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordComment:     false,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	var comment string
	if (*existMap)[args.FieldWordComment] {
		comment = (*valueMap)[args.FieldWordComment].(string)
	} else {
		comment = ""
	}
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// gen ak & sk
	var accessKey, secretKey string
	accessKey = code.GenAccessKey()
	secretKey = code.GenSecretKey()

	// save it into mongodb 由于 scheduler 尚未实现 key 的同步， 在本地实现 ak sk 添加
	//insertKeySuccess := dao.AccessKeyDao.InsertKey(userID, accessKey, secretKey, comment)
	//if !checkDaoSuccess(con, insertKeySuccess) {
	//	return
	//}

	//sync it with scheduler 由于 scheduler 尚未实现 key 的同步，先禁用该同步
	key := model.AccessKey{
		UserID:     userID,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Comment:    comment,
		Available:  true,
		CreateTime: time.Now(),
	}
	postKeyResponse, postKeySuccess := scheduler.PostKeyToScheduler(&key)
	if !postKeySuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postKeyResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": postKeyResponse.Code,
			"msg":  postKeyResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// success with gen keys
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "生成密钥成功",
		"data": gin.H{
			"AccessKey": accessKey,
			"SecretKey": secretKey,
		},
	})
}

func UserDeleteKey(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordAccessKey:   true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	accessKey := (*valueMap)[args.FieldWordAccessKey].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}

	// delete keys with dao
	//deleteKeyResult, deleteKeySuccess := dao.AccessKeyDao.DeleteKey(userID, accessKey)
	//if !checkDaoSuccess(con, deleteKeySuccess) {
	//	return
	//}
	// no key has been deleted
	//if deleteKeyResult.DeletedCount == 0 {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDeleteNothing,
	//		"msg":  "此密钥本来就不存在,你删nmn¿",
	//		"data": gin.H{},
	//	})
	//	return
	//}

	// delete & sync with scheduler
	deleteKeyResponse, deleteKeySuccess := scheduler.DeleteKeyToScheduler(userID, accessKey)
	if !deleteKeySuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if deleteKeyResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": deleteKeyResponse.Code,
			"msg":  deleteKeyResponse.Msg,
			"data": gin.H{},
		})
		return
	}
	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户删除密钥成功",
		"data": gin.H{},
	})
}

func UserChangeKeyStatus(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordAccessKey:   true,
		args.FieldWordStatus:      true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	accessKey := (*valueMap)[args.FieldWordAccessKey].(string)
	status := (*valueMap)[args.FieldWordStatus].(bool)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}

	// change key status with dao
	//changeKeyStatusResult, changeKeyStatusSuccess := dao.AccessKeyDao.ChangeKeyStatus(userID, accessKey, status)
	//if !checkDaoSuccess(con, changeKeyStatusSuccess) {
	//	return
	//}
	// no key has been changed
	//if changeKeyStatusResult.MatchedCount == 0 {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDeleteNothing,
	//		"msg":  "此密钥不存在,你改nmn¿",
	//		"data": gin.H{},
	//	})
	//	return
	//}

	// get origin key
	getKeyResult, getKeySuccess := dao.AccessKeyDao.GetWithAccessKey(userID, accessKey)
	if !checkDaoSuccess(con, getKeySuccess) {
		return
	}
	getKeyResult.Available = status
	// sync with scheduler
	postKeyResponse, postKeySuccess := scheduler.PostKeyToScheduler(getKeyResult)
	if !postKeySuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postKeyResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": postKeyResponse.Code,
			"msg":  postKeyResponse.Msg,
			"data": gin.H{},
		})
		return
	}

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户修改密钥状态成功",
		"data": gin.H{},
	})
}

func UserChangeKeyComment(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordAccessKey:   true,
		args.FieldWordComment:     false,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	accessKey := (*valueMap)[args.FieldWordAccessKey].(string)
	var newComment string
	if (*existMap)[args.FieldWordComment] {
		newComment = (*valueMap)[args.FieldWordComment].(string)
	} else {
		newComment = ""
	}
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}

	// change key comment with dao
	//changeKeyCommentResult, changeKeyCommentSuccess := dao.AccessKeyDao.ChangeKeyComment(userID, accessKey, newComment)
	//if !checkDaoSuccess(con, changeKeyCommentSuccess) {
	//	return
	//}
	// no key has been changed
	//if changeKeyCommentResult.MatchedCount == 0 {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDeleteNothing,
	//		"msg":  "此密钥不存在,你改nmn¿",
	//		"data": gin.H{},
	//	})
	//	return
	//}

	// get origin key
	getKeyResult, getKeySuccess := dao.AccessKeyDao.GetWithAccessKey(userID, accessKey)
	if !checkDaoSuccess(con, getKeySuccess) {
		return
	}
	getKeyResult.Comment = newComment
	// sync with scheduler
	postKeyResponse, postKeySuccess := scheduler.PostKeyToScheduler(getKeyResult)
	if !postKeySuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postKeyResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": postKeyResponse.Code,
			"msg":  postKeyResponse.Msg,
			"data": gin.H{},
		})
		return
	}

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户修改密钥备注成功",
		"data": gin.H{},
	})
}

func UserRemakeKey(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordAccessKey:   true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	accessKey := (*valueMap)[args.FieldWordAccessKey].(string)
	// check token
	userID, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	if !valid {
		return
	}
	// generate a new secret key
	var secretKey string
	secretKey = code.GenSecretKey()

	// remake secret key with dao
	//remakeKeyResult, remakeKeySuccess := dao.AccessKeyDao.RemakeKey(userID, accessKey, secretKey)
	//if !checkDaoSuccess(con, remakeKeySuccess) {
	//	return
	//}
	// no key has been remade
	//if remakeKeyResult.MatchedCount == 0 {
	//	con.JSON(http.StatusOK, gin.H{
	//		"code": args.CodeDeleteNothing,
	//		"msg":  "此密钥不存在,你/remake nmn¿",
	//		"data": gin.H{},
	//	})
	//	return
	//}

	// get origin key
	getKeyResult, getKeySuccess := dao.AccessKeyDao.GetWithAccessKey(userID, accessKey)
	if !checkDaoSuccess(con, getKeySuccess) {
		return
	}
	getKeyResult.SecretKey = secretKey
	// sync with scheduler
	postKeyResponse, postKeySuccess := scheduler.PostKeyToScheduler(getKeyResult)
	if !postKeySuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeJsonError,
			"msg":  "解析scheduler-json信息有误",
			"data": gin.H{},
		})
		return
	}
	if postKeyResponse.Code != args.CodeOK {
		// error in scheduler
		con.JSON(http.StatusOK, gin.H{
			"code": postKeyResponse.Code,
			"msg":  postKeyResponse.Msg,
			"data": gin.H{},
		})
		return
	}

	// success
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "用户重置密钥成功",
		"data": gin.H{},
	})
}

//func UserUploadAvatar(c *gin.Context) {
//	userID := context.PostForm("user_id")
//	file, err := context.FormFile("avatar")
//	if err!= nil {
//		tool.Failed(context, "参数解析失败")
//		return
//	}
//
//	num := "user_"+userID
//	fmt.Println("num:",num)
//	//2. 只有登录才能修改用户头像信息
//	//sesstion := tool.Getsess(context, num)
//	//fmt.Println("sesstion:",sesstion)
//	//if sesstion == nil {
//	//	tool.Failed(context, "参数不合法")
//	//	return
//	//}
//	var member model.Member = model.Member{
//		Id:           1,
//	}
//	//json.Unmarshal(sesstion.([]byte),&member)
//
//	//2. file保存在本地
//	fileName := "./uploadfile/"+strconv.FormatInt(time.Now().Unix(), 10)+file.Filename
//	err = context.SaveUploadedFile(file, fileName)
//	if err != nil {
//		tool.Failed(context,"头像更新失败")
//		return
//	}
//
//	//将文件上传到fastDFS系统
//	fileId := tool.UploadFile(fileName)
//	if fileId != "" {
//		//删除本地文件
//		os.Remove(fileName)
//		//3. 将保存后的文件本地路径, 保存到用户表中的头像字段
//		memberService := service.MemberService{}
//		fileN := memberService.UploadAvator(member.Id, fileId)
//		if len(fileN) != 0 {
//			tool.Success(context, tool.FileServerAddr()+"/" +fileN)
//			return
//		}
//	}
//	//4. 返回结果
//	tool.Failed(context, "上传失败")
//
//}
