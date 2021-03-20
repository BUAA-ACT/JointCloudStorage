package controller

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func CheckErr(err error, label string) bool {
	if err != nil {
		logrus.Warnf("%v ERR: %v", label, err)
		return true
	}
	return false
}

const (
	SECRET = "develop" // todo 签名密码加入配置文件
)

type FileAccessClaims struct {
	path string
	uid  string
	jwt.StandardClaims
}

func GenerateLocalFileAccessToken(path string, uid string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, FileAccessClaims{
		path: path,
		uid:  uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			Issuer:    "transporter",
		},
	})
	rs, err := token.SignedString([]byte(SECRET))
	return rs, err
}

func ParseLocalFileAccessToken(accessToken string) (*FileAccessClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &FileAccessClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*FileAccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func genRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	)
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken := c.DefaultQuery("token", "")
		if accessToken == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "no token found",
			})
			c.Abort()
			return
		}
		mc, err := ParseLocalFileAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("filePath", mc.path)
		c.Next() // 后续的处理函数可以用过c.Get("filePath")来获取当前请求的用户信息
	}
}
