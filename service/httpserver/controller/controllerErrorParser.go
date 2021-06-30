package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/regex"
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/textproto"
	"strings"
)

func getValueAndExist(con *gin.Context, fields *map[string]bool) (*map[string]interface{}, *map[string]bool) {
	valueMap := make(map[string]interface{})
	existMap := make(map[string]bool)
	if con.Request.Method == args.HttpMethodGet {
		for field := range *fields {
			getValue, ok := con.GetQuery(field)
			valueMap[field] = getValue
			existMap[field] = ok
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
				existMap[field] = ok
			}
		} else if strings.Contains(httpType, args.HttpContentTypeUrlEncoded) {
			for field := range *fields {
				encodeValue, ok := con.GetPostForm(field)
				valueMap[field] = encodeValue
				existMap[field] = ok
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
				fmt.Println("error in get value:")
				fmt.Println(err)
			}
			for field := range *fields {
				jsonValue, ok := result[field]
				valueMap[field] = jsonValue
				existMap[field] = ok
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
			existMap[field] = err == nil
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
				valueMap[field] = header[newField][0]
				existMap[field] = true
			}
		}
	}
	return &valueMap, &existMap
}

func getQueryAndReturn(con *gin.Context, fields *map[string]bool) (*map[string]interface{}, *map[string]bool) {
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

func UserCheckAccessToken(con *gin.Context, accessToken string, permitRoles *[]string) (string, bool) {
	userID, valid := dao.AccessTokenDao.CheckValid(accessToken)
	if !valid {
		con.JSON(http.StatusOK, gin.H{
			"code": args.CodeInvalidAccessToken,
			"msg":  "用户令牌无效",
			"data": gin.H{},
		})
		return "", false
	}
	user, valid := dao.UserDao.GetUserInfo(userID)
	for _, role := range *permitRoles {
		// role == UserAllRole can be all permit
		if user.Role == role || role == args.UserAllRole {
			return userID, true
		}
	}
	con.JSON(http.StatusOK, gin.H{
		"code": args.CodeInvalidRole,
		"msg":  "用户权限错误",
		"data": gin.H{},
	})
	return "", false
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
