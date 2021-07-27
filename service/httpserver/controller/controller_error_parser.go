package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/regex"
	"cloud-storage-httpserver/service/tools"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/textproto"
	"strings"
)

func putValueIntoMap(value interface{}, ok bool, field string, valueMap *map[string]interface{}, existMap *map[string]bool) {
	if value == nil {
		(*valueMap)[field] = ""
		(*existMap)[field] = false
	} else {
		(*valueMap)[field] = value
		(*existMap)[field] = ok
	}
	return
}

func getValueAndExist(con *gin.Context, fields *map[string]bool) (*map[string]interface{}, *map[string]bool) {
	valueMap := make(map[string]interface{})
	existMap := make(map[string]bool)
	if con.Request.Method == args.HttpMethodGet {
		for field := range *fields {
			getValue, ok := con.GetQuery(field)
			valueMap[field] = getValue
			existMap[field] = ok && getValue != ""
		}
	} else if con.Request.Method == args.HttpMethodPost {
		httpType := con.GetHeader("Content-Type")
		if strings.Contains(httpType, args.HttpContentTypeFormData) {
			var ok bool
			var data interface{}
			var value string
			for field := range *fields {
				value, ok = con.GetPostForm(field)
				if !ok {
					file, fileErr := con.FormFile(field)
					data = file
					ok = fileErr == nil
				} else {
					data = value
				}
				valueMap[field] = data
				existMap[field] = ok && data != nil
			}
		} else if strings.Contains(httpType, args.HttpContentTypeUrlEncoded) {
			for field := range *fields {
				encodeValue, ok := con.GetPostForm(field)
				valueMap[field] = encodeValue
				existMap[field] = ok && encodeValue != ""
			}
		} else if strings.Contains(httpType, args.HttpContentTypeRaw) {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		} else if strings.Contains(httpType, args.HttpContentTypeHTML) {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		} else if strings.Contains(httpType, args.HttpContentTypeJson) {
			var result map[string]interface{}
			err := con.ShouldBindJSON(&result)
			//resultJson,err := json.Marshal(result)
			//con.Request.Body = ioutil.NopCloser(bytes.NewBuffer(resultJson))
			if err != nil {
				log.Println("error in get value:")
				log.Println(err)
			}
			for field := range *fields {
				jsonValue, ok := result[field]
				valueMap[field] = jsonValue
				existMap[field] = ok && jsonValue != nil
			}
		} else if strings.Contains(httpType, args.HttpContentTypeJavascript) {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		} else if strings.Contains(httpType, args.HttpContentTypeXML) {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		} else if strings.Contains(httpType, args.HttpContentTypeMS) {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		} else {
			for field := range *fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		}
	} else {
		for field := range *fields {
			valueMap[field] = ""
			existMap[field] = false
		}
	}
	// get field from cookies
	for field := range *fields {
		if !existMap[field] {
			cookieValue, err := con.Cookie(field)
			valueMap[field] = cookieValue
			existMap[field] = (err == nil) && (cookieValue != "")
		}
	}
	// get field from headers
	for field := range *fields {
		// pass cookies
		if field == "Cookie" {
			continue
		}
		if !existMap[field] {
			// get it with origin form
			newField := textproto.CanonicalMIMEHeaderKey(field)
			header := con.Request.Header
			if header[newField] != nil && len(header[newField]) != 0 {
				headerValue := header[newField][0]
				valueMap[field] = headerValue
				existMap[field] = headerValue != ""
			}
		}
	}
	return &valueMap, &existMap
}

func getQueryAndReturnWithHttp(con *gin.Context, fields *map[string]bool) (*map[string]interface{}, *map[string]bool) {
	fieldValues, fieldExists := getValueAndExist(con, fields)
	for field, fieldExist := range *fieldExists {
		if !fieldExist && (*fields)[field] {
			con.JSON(http.StatusOK, gin.H{
				"code": args.CodeFieldNotExist,
				"msg":  "没有" + field + "字段",
				"data": gin.H{},
			})
			return fieldValues, fieldExists
		}
		realValue, regexSuccess := regex.CheckRegex((*fieldValues)[field], field)
		if !regexSuccess {
			con.JSON(http.StatusOK, gin.H{
				"code": args.CodeRegexWrong,
				"msg":  field + "字段格式错误",
				"data": gin.H{},
			})
			return fieldValues, fieldExists
		}
		(*fieldExists)[field] = fieldExist && regexSuccess
		(*fieldValues)[field] = realValue
	}
	return fieldValues, fieldExists
}

func getVerifyAndReturnWithWebSocket(ws *websocket.Conn, fields *map[string]bool) (*map[string]interface{}, *map[string]bool) {
	// get json data from ws
	var jsonMap map[string]interface{}
	jsonErr := ws.ReadJSON(&jsonMap)
	if jsonErr != nil {
		log.Println("fucking reading json problem while json error in head: " + jsonErr.Error())
		returnMap := gin.H{
			"code": args.CodeJsonError,
			"msg":  "json解析有误",
			"data": gin.H{},
		}
		writeJsonErr := ws.WriteJSON(returnMap)
		if writeJsonErr != nil {
			tools.PrintError(writeJsonErr)
		}
		return nil, nil
	}
	// get field value and exist
	fieldValues := make(map[string]interface{})
	fieldExists := make(map[string]bool)
	for field := range *fields {
		jsonValue, jsonValueOk := jsonMap[field]
		fieldValues[field] = jsonValue
		fieldExists[field] = jsonValueOk && jsonValue != nil
	}

	for field, fieldExist := range fieldExists {
		if !fieldExist && (*fields)[field] {
			returnMap := gin.H{
				"code": args.CodeFieldNotExist,
				"msg":  "没有" + field + "字段",
				"data": gin.H{},
			}
			writeJsonErr := ws.WriteJSON(returnMap)
			if writeJsonErr != nil {
				tools.PrintError(writeJsonErr)
			}
			return &fieldValues, &fieldExists
		}
		realValue, regexSuccess := regex.CheckRegex(fieldValues[field], field)
		if !regexSuccess {
			returnMap := gin.H{
				"code": args.CodeRegexWrong,
				"msg":  field + "字段格式错误",
				"data": gin.H{},
			}
			writeJsonErr := ws.WriteJSON(returnMap)
			if writeJsonErr != nil {
				tools.PrintError(writeJsonErr)
			}
			return &fieldValues, &fieldExists
		}
		fieldExists[field] = fieldExist && regexSuccess
		fieldValues[field] = realValue
	}
	return &fieldValues, &fieldExists
}

func UserCheckAccessToken(con *gin.Context, accessToken string, permitRoles *[]string) (*model.User, bool) {
	userID, valid := dao.AccessTokenDao.CheckValid(accessToken)
	if !valid {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeInvalidAccessToken,
			"msg":  "用户令牌无效",
			"data": gin.H{},
		})
		return nil, false
	}
	user, valid := dao.UserDao.GetUserInfo(userID)
	for _, role := range *permitRoles {
		// role == UserAllRole can be all permit
		if user.Role == role || role == args.UserAllRole {
			con.SetCookie("AccessToken", accessToken, int(args.AccessTokenTimeOut.Seconds()), "/", strings.Split(con.Request.Host, ":")[0], false, false)
			return user, true
		}
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeInvalidRole,
		"msg":  "用户权限错误",
		"data": gin.H{},
	})
	return nil, false
}

func UserCheckStatus(con *gin.Context, user *model.User, statusMap *map[string]bool) bool {
	var code int
	var message string
	for field := range *statusMap {
		if user.Status == field {
			code, message = tools.UserStatusMessageCode(field)
			con.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  message,
				"data": gin.H{},
			})
			return false
		}
	}
	return true
}

func checkDaoSuccess(con *gin.Context, success bool) bool {
	if success {
		return true
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeDatabaseError,
		"msg":  "数据库错误",
		"data": gin.H{},
	})
	return false
}
