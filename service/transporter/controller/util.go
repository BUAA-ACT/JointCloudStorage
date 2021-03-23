package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"context"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

type AuthClaims struct {
	Path string
	Uid  string
	Tid  string
	jwt.StandardClaims
}

func GenerateTaskAccessToken(tid string, uid string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		Uid: uid,
		Tid: tid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			Issuer:    "transporter",
		},
	})
	rs, err := token.SignedString([]byte(SECRET))
	return rs, err
}

func GenerateLocalFileAccessToken(path string, uid string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		Path: path,
		Uid:  uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			Issuer:    "transporter",
		},
	})
	rs, err := token.SignedString([]byte(SECRET))
	return rs, err
}

func ParseLocalFileAccessToken(accessToken string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
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
			accessToken = c.DefaultPostForm("token", "")
		}
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
		c.Set("filePath", mc.Path)
		c.Set("tokenUid", mc.Uid)
		c.Set("tokenTid", mc.Tid)
		c.Next() // 后续的处理函数可以用过c.Get("filePath")来获取当前请求的用户信息
	}
}

// 获取文件 ContentType
func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func FromFileInfoGetUidAndPath(file *model.File) (uid string, path string) {
	p := strings.Index(file.Id, "/")
	uid = file.Id[0:p]
	path = file.Id[p+1:]
	return
}

func ClearAll() {
	clientOptions := options.Client().ApplyURI("mongodb://192.168.105.8:20100")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	defer client.Disconnect(context.TODO())
	if err != nil {
		log.Print(err)
	}
	collection := client.Database("transporterTasks").Collection("Tasks")
	collection.Drop(context.TODO())
	collection = client.Database("Cloud").Collection("FileDatabase")
	collection.Drop(context.TODO())
	return
}
