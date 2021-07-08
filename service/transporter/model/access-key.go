package model

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

// AccessKey 统一编程接口用户认证秘钥
type AccessKey struct {
	UserId     string    `bson:"user_id"`
	AccessKey  string    `bson:"access_key"`
	SecretKey  string    `bson:"secret_key"`
	CreateTime time.Time `bson:"create_time"`
	Available  bool      `bson:"available"`
}

// AccessKeyDB 用户秘钥数据库
type AccessKeyDB struct {
	Dao *Dao
}

// NewAccessKeyDB 用于初始化用户秘钥数据库连接
// Dao mongodb 数据库连接 dao
func NewAccessKeyDB(dao *Dao) (db *AccessKeyDB, err error) {
	if dao.TestPing() {
		db = &AccessKeyDB{Dao: dao}
		return db, nil
	} else {
		return nil, errors.New(util.ErrorDBConnection)
	}
}

func GenNewRandom32str() string {
	s := uuid.New().String()
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// GenerateKeys 为用户生成一套秘钥，开发测试使用，生产环境该环节由 httpserver 完成
func (db *AccessKeyDB) GenerateKeys(userId string) (AccessKey, error) {
	col := db.Dao.GetCollection(AccessKeyCollection)
	ak := GenNewRandom32str()
	sk := GenNewRandom32str()

	accessKey := AccessKey{
		UserId:     userId,
		AccessKey:  ak,
		SecretKey:  sk,
		CreateTime: time.Now(),
		Available:  true,
	}
	_, err := col.InsertOne(context.TODO(), accessKey)
	return accessKey, err
}

func (db *AccessKeyDB) Certificate(ak string, sk string) (userId string, err error) {
	col := db.Dao.GetCollection(AccessKeyCollection)
	var accessKey AccessKey
	err = col.FindOne(context.TODO(), bson.M{"access_key": ak}).Decode(&accessKey)
	if err != nil {
		return "", err
	}
	if accessKey.SecretKey != sk {
		return "", errors.New(util.ErrorCertificate)
	}
	return accessKey.UserId, nil
}

func (db *AccessKeyDB) GetKey(ak string) (*AccessKey, error) {
	col := db.Dao.GetCollection(AccessKeyCollection)
	var accessKey AccessKey
	err := col.FindOne(context.TODO(), bson.M{"access_key": ak}).Decode(&accessKey)
	if err != nil {
		return nil, err
	} else {
		return &accessKey, nil
	}
}
