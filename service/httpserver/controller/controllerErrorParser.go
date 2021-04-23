package controller

import (
	. "cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/service/regex"
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func getValueAndExist(con *gin.Context, fields map[string]bool) (map[string]interface{}, map[string]bool) {
	valueMap := make(map[string]interface{})
	existMap := make(map[string]bool)
	if con.Request.Method == HttpMethodGet {
		for field := range fields {
			getValue, ok := con.GetQuery(field)
			valueMap[field] = getValue
			existMap[field] = ok
		}
		return valueMap, existMap
	} else if con.Request.Method == HttpMethodPost {
		httpType := con.GetHeader("Content-Type")
		if strings.Contains(httpType, HttpContentTypeUrlEncoded) {
			for field := range fields {
				encodeValue, ok := con.GetPostForm(field)
				valueMap[field] = encodeValue
				existMap[field] = ok
			}
		} else if strings.Contains(httpType, HttpContentTypeJson) {
			var result map[string]interface{}
			err := con.ShouldBindJSON(&result)
			//resultJson,err := json.Marshal(result)
			//con.Request.Body = ioutil.NopCloser(bytes.NewBuffer(resultJson))
			if err != nil {
				fmt.Println("error in get value:")
				fmt.Println(err)
			}
			for field := range fields {
				jsonValue, ok := result[field]
				valueMap[field] = jsonValue
				existMap[field] = ok
			}
		} else {
			for field := range fields {
				valueMap[field] = ""
				existMap[field] = false
			}
		}
	} else {
		for field := range fields {
			valueMap[field] = ""
			existMap[field] = false
		}
	}
	for field := range fields {
		if !existMap[field] {
			cookieValue, err := con.Cookie(field)
			ok := tools.PrintError(err)
			valueMap[field] = cookieValue
			existMap[field] = ok
		}
	}
	return valueMap, existMap
}

func getQueryAndReturn(con *gin.Context, fields map[string]bool) (map[string]interface{}, map[string]bool) {
	fieldValues, fieldExists := getValueAndExist(con, fields)
	for field, fieldExist := range fieldExists {
		if !fieldExist && fields[field] {
			con.JSON(http.StatusOK, gin.H{
				"code": CodeFieldNotExist,
				"msg":  "没有" + field + "字段",
				"data": gin.H{},
			})
			return fieldValues, fieldExists
		}
		realValue, regexSuccess := regex.CheckRegex(fieldValues[field], field)
		if !regexSuccess {
			con.JSON(http.StatusOK, gin.H{
				"code": CodeRegexWrong,
				"msg":  field + "字段格式错误",
				"data": gin.H{},
			})
			return fieldValues, fieldExists
		}
		fieldExists[field] = fieldExist && regexSuccess
		fieldValues[field] = realValue
	}
	return fieldValues, fieldExists
}

func UserCheckAccessToken(con *gin.Context, accessToken string) (string, bool) {
	userId, valid := dao.AccessTokenDao.CheckValid(accessToken)
	if !valid {
		con.JSON(http.StatusOK, gin.H{
			"code": CodeTokenNotValid,
			"msg":  "用户令牌无效",
			"data": gin.H{},
		})
		return "", false
	}
	return userId, true
}

func checkDaoError(con *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	con.JSON(http.StatusOK, gin.H{
		"code": CodeDatabaseError,
		"msg":  "数据库错误",
		"data": gin.H{},
	})
	return true
}
