package storageInterface

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// method
	HttpMethodGet    = "GET"
	HttpMethodPost   = "POST"
	HttpMethodPut    = "PUT"
	HttpMethodDelete = "DELETE"

	// content type
	HttpContentTypeUrlEncoded = "application/x-www-form-urlencoded"
	HttpContentTypeFormData   = "multipart/form-data"
	HttpContentTypeRaw        = "text/plain"
	HttpContentTypeHTML       = "text/html"
	HttpContentTypeJavascript = "application/javascript"
	HttpContentTypeJson       = "application/json"
	HttpContentTypeXML        = "application/xml"
	HttpContentTypeMS         = "application/x-msdownload"

	// key
	HttpHeaderKeyAuthorization = "Authorization"
	HttpHeaderKeyTime          = "Time"
	HttpHeaderMethod           = "Method"
	HttpHeaderURL              = "URL"

	// length
	AKLENGTH        = 32
	SKLENGTH        = 32
	DIGESTLENGTH    = 48
	SIGNATURELENGTH = 64
	JCSMARK         = "JCS"
)

// JSIAuthMiddleware  JSI 的认证中间件
func (jsi *JointStorageInterface) JSIAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//c.JSON(http.StatusUnauthorized, "")
		//c.Abort()
		uid, rs := jsi.auth(c)
		if !rs {
			c.JSON(http.StatusUnauthorized, "")
			c.Abort()
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("uid", uid)
		userInfo, err := jsi.processor.UserDatabase.GetUserFromID(uid)
		if err != nil {
			util.Log(logrus.ErrorLevel, "JSI GetObjectList", "get Userinfo fail",
				"", "err", err.Error())
			c.String(http.StatusInternalServerError, "")
			c.Abort()
		}
		c.Set("userInfo", userInfo)
		c.Next() // 后续的处理函数可以用过c.Get("filePath")来获取当前请求的用户信息
	}
}

func (jsi *JointStorageInterface) auth(con *gin.Context) (uid string, result bool) {
	// pure what we need into key
	valueMap, err := PureIntoKey(con)
	if err != nil {
		con.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code": 110,
			"msg":  "您的时间好像被虵偷走了",
			"data": gin.H{},
		})
		return "", false
	}
	// check authorizarion
	accessKey, requestEncodeSign, getAuthSuccess := CheckAuthorization(con)
	if !getAuthSuccess {
		con.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code": 111,
			"msg":  "认证字段不存在",
			"data": gin.H{},
		})
		return "", false
	}
	// get all keys and check it
	keys := []string{HttpHeaderMethod, HttpHeaderURL, HttpHeaderKeyTime}
	// generate origin sign
	var sign string = ""
	for _, key := range keys {
		value, exist := valueMap[key]
		if !exist {
			con.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 111,
				"msg":  key + "字段不存在",
				"data": gin.H{},
			})
			return "", false
		}
		sign = sign + key + ":" + value + "\r\n"
	}
	key, err := jsi.processor.AccessKeyDatabase.GetKey(accessKey)
	if err != nil {
		return "", false
	}
	// compute with sha3-256
	signature, signErr := Sha3Encode(sign, key.SecretKey)
	if signErr != nil {
		// error with sk
		con.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code": 112,
			"msg":  signErr.Error(),
			"data": gin.H{},
		})
		return "", false
	}
	requestSignature, err := base64.URLEncoding.DecodeString(requestEncodeSign)
	if err != nil {
		return "", false
	}
	if signature != string(requestSignature) {
		// error with the code
		con.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code": 113,
			"msg":  "验证密码错误",
			"data": gin.H{},
		})
		return "", false
	}
	// permit its sk
	return key.UserId, true
}

func isSpace(ch byte) bool {
	return ch == byte(' ') || ch == byte('\t') || ch == byte('\r') || ch == byte('\n')
}

func isHexDig(ch byte) bool {
	return (byte('0') <= ch && ch <= byte('9')) || (byte('a') <= ch && ch <= byte('f')) || (byte('A') <= ch && ch <= byte('F'))
}

func isURLBase64Encode(ch byte) bool {
	return (byte('0') <= ch && ch <= byte('9')) || (byte('a') <= ch && ch <= byte('z')) || (byte('A') <= ch && ch <= byte('Z') || ch == byte('-') || ch == byte('_'))

}

func CheckAuthorization(con *gin.Context) (string, string, bool) {
	authHead := con.Request.Header.Get(HttpHeaderKeyAuthorization)
	if authHead == "" {
		return "", "", false
	}
	var prefix strings.Builder
	var accessKey strings.Builder
	var signature strings.Builder
	var status = 0
	for i := 0; i < len(authHead); i++ {
		if isSpace(authHead[i]) {
			continue
		} else if status == 0 {
			if authHead[i] == byte('(') {
				// start with JCS
				if prefix.String() == JCSMARK {
					status = 1
				} else {
					return "", "", false
				}
			} else {
				prefix.WriteByte(authHead[i])
			}
		} else if status == 1 {
			// ':' cut the AK and Signature
			if authHead[i] == byte(':') {
				// AK is 32 byte
				if accessKey.Len() != AKLENGTH {
					return "", "", false
				}
				status = 2
			} else if isHexDig(authHead[i]) {
				accessKey.WriteByte(authHead[i])
			}
		} else if status == 2 {
			if authHead[i] == byte(')') {
				// signature is 64 byte
				if signature.Len() != SIGNATURELENGTH {
					return "", "", false
				}
				status = 3
			} else if isURLBase64Encode(authHead[i]) {
				signature.WriteByte(authHead[i])
			}
		} else {
			return "", "", false
		}
	}
	if status != 3 {
		return "", "", false
	}
	return accessKey.String(), signature.String(), true
}

func PureIntoKey(con *gin.Context) (keys map[string]string, err error) {
	requestTime := con.Request.Header.Get(HttpHeaderKeyTime)
	if requestTime == "" {
		return nil, errors.New("can't get time")
	}
	// use Unix second
	realTime, timeParseErr := strconv.ParseInt(requestTime, 10, 64)
	if timeParseErr != nil {
		return nil, errors.New("can't get time")
	}
	nowTime := time.Now().Unix()
	if realTime > nowTime+60 || realTime+15*60 < nowTime {
		return nil, errors.New("invalid time")
	}
	method := con.Request.Method
	url := con.Request.URL.String()
	url = strings.TrimSuffix(url, "/")
	// pure header into key
	valueMap := map[string]string{
		HttpHeaderMethod:  method,
		HttpHeaderURL:     url,
		HttpHeaderKeyTime: strconv.FormatInt(realTime, 10),
	}
	return valueMap, nil
}

func JSISign(r *http.Request, ak string, sk string) (rs *http.Request, err error) {
	wholeURL := strings.TrimSuffix(r.URL.String(), "/")
	method := r.Method

	//path := r.TLS
	// use Unix second
	nowTime := strconv.FormatInt(time.Now().Unix(), 10)
	// generate the origin sign
	keys := []string{HttpHeaderMethod, HttpHeaderURL, HttpHeaderKeyTime}
	valueMap := map[string]string{
		HttpHeaderMethod:  method,
		HttpHeaderURL:     wholeURL,
		HttpHeaderKeyTime: nowTime,
	}
	var sign string = ""
	for _, key := range keys {
		sign = sign + key + ":" + valueMap[key] + "\r\n"
	}
	// compute with sha3-256
	signature, signErr := Sha3Encode(sign, sk)
	if signErr != nil {
		return nil, signErr
	}
	// set in base64 url encode with padding
	encodeSign := base64.URLEncoding.EncodeToString([]byte(signature))
	// out with format
	authorization := "JCS(" + ak + ":" + encodeSign + ")"
	// send with method,url,accessKey,authorization
	r.Header.Set(HttpHeaderKeyTime, nowTime)
	r.Header.Set(HttpHeaderKeyAuthorization, authorization)
	return r, nil
}
