package controller

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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
	token := jwt.NewWithClaims(jwt.SigningMethodES256, FileAccessClaims{
		path: path,
		uid:  uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			Issuer:    "transporter",
		},
	})
	return token.SignedString(SECRET)
}
