package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	email := (*valueMap)[args.FieldWordEmail].(string)
	password := (*valueMap)[args.FieldWordPassword].(string)
	nickname := (*valueMap)[args.FieldWordNickname].(string)
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
	userId := email
	// save with dao
	user := &model.User{
		UserId:       userId,
		Email:        email,
		Password:     code.AesEncrypt(password, *args.EncryptKey),
		Nickname:     nickname,
		Role:         args.UserHostRole,
		Avatar:       "default-avatar.png",
		CreateTime:   nowTime,
		LastModified: nowTime,
		Status:       args.UserNormalStatus,
	}
	dao.UserDao.CreateNewUser(*user)
	// record verify code
	verifyCode := code.GenVerifyCode()
	dao.VerifyCodeDao.InsertVerifyCode(email, verifyCode)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
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

	success := dao.VerifyCodeDao.VerifyEmail(email, verifyCode)
	if success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeOK,
			"msg":  "验证码正确",
			"data": gin.H{},
		})
		// update user status
		dao.UserDao.SetUserStatusWithEmail(email, args.UserNormalStatus)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
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
	dao.AccessTokenDao.InsertAccessToken(token, user.UserId)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	// check token is valid
	_, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	// delete token
	dao.AccessTokenDao.DeleteAccessToken(accessToken)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	//check token
	_, valid := UserCheckAccessToken(con, accessToken)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	originPassword := (*valueMap)[args.FieldWordOriginPassword].(string)
	newPassword := (*valueMap)[args.FieldWordNewPassword].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	// verify role
	user, success := dao.UserDao.GetUserInfo(userId)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	if user.Role == args.UserGuestRole {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "Guest禁止修改密码",
			"data": gin.H{},
		})
		return
	}

	// verify origin password
	loginSuccess := dao.UserDao.LoginWithId(userId, originPassword)
	if !loginSuccess {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodePasswordNotRight,
			"msg":  "原密码错误",
			"data": gin.H{},
		})
		return
	}
	// dao change password
	dao.UserDao.SetUserPassword(userId, newPassword)
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "修改密码成功",
		"data": gin.H{},
	})
}

func UserChangeEmail(con *gin.Context) {
	fieldRequired := map[string]bool{
		args.FieldWordAccessToken: true,
		args.FieldWordNewEmail:    true,
	}
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	newEmail := (*valueMap)[args.FieldWordNewEmail].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	// can't change email now!
	user, success := dao.UserDao.GetUserInfo(userId)
	if !success {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeDatabaseError,
			"msg":  "数据库错误",
			"data": gin.H{},
		})
		return
	}
	if user.Role == args.UserHostRole || user.Role == args.UserGuestRole {
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
	dao.UserDao.SetUserEmail(userId, newEmail)
	dao.UserDao.SetUserStatusWithId(userId, args.UserVerifyStatus)
	// verify code
	verifyCode := code.GenVerifyCode()
	dao.VerifyCodeDao.InsertVerifyCode(newEmail, verifyCode)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	newNickname := (*valueMap)[args.FieldWordNickname].(string)
	// check token
	userId, valid := UserCheckAccessToken(con, accessToken)
	if !valid {
		return
	}
	dao.UserDao.SetUserNickname(userId, newNickname)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
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
	valueMap, existMap := getQueryAndReturn(con, &fieldRequired)
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
	userId, valid := UserCheckAccessToken(con, accessToken)
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
	dao.UserDao.SetUserPreference(userId, preference)
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeOK,
		"msg":  "设置个人偏好成功",
		"data": gin.H{},
	})
}

//func UserUploadAvatar(c *gin.Context) {
//	userId := context.PostForm("user_id")
//	file, err := context.FormFile("avatar")
//	if err!= nil {
//		tool.Failed(context, "参数解析失败")
//		return
//	}
//
//	num := "user_"+userId
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
