package config

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
)

// Version constant */
const (
	Version = "v0.2"
)

/* http constant */
const (
	CloudCollectionName           = "Cloud"
	TempCloudCollectionName       = "TempCloud"
	VoteCloudCollectionName       = "VoteCloud"
	UserCollectionName            = "User"
	FileCollectionName            = "File"
	MigrationAdviceCollectionName = "MigrationAdvice"
	AccessKeyCollectionName       = "AccessKey"

	codeOK            = 200
	codeBadRequest    = 400
	codeUnauthorized  = 401
	codeInternalError = 500
)

/* user role constant */
const (
	RoleHost  = "HOST"
	RoleGuest = "GUEST"
)

/* Advice status */
const (
	AdviceStatusPending = "PENDING"

	AdviceStatus = "TODO"
)

/* */
var (
	FlagMongo              = flag.String("mongo", "mongodb://localhost:27017", "mongodb address")
	FlagAddress            = flag.String("addr", ":8082", "scheduler address")
	FlagEnv                = flag.String("env", "test", "dev|test|prod")
	FlagCloudID            = flag.String("cid", "aliyun-beijing", "cloud id")
	FlagAESKey             = flag.String("aes", "1234567890123456", "aes key")
	FlagRescheduleInterval = flag.Duration("reschedule", time.Minute*1, "reschedule interval")
	FlagHeartbeatInterval  = flag.Duration("heartbeat", time.Second*30, "heartbeat interval")
	AddrMap                = make(map[string]string)
)

var conf *Config
var lock sync.Mutex
var lock sync.RWMutex

func RegisterDao(URI string, client ClientConfig) {

	for databaseName, database := range client.Databases {

	}
}

func SetGlobalConfig() {
	lock.Lock()
	conf = Factory()
	a := reflect.ValueOf(conf).Elem().String()
	fmt.Println(a)

}

func GetConfig() *Config {
	lock.Lock()
	if conf == nil {
		logrus.Errorf("读取配置文件时，config 尚未设置！！正在自动配置中")
		conf = &Config{
			FlagMongo:              *FlagMongo,
			FlagAddress:            *FlagAddress,
			FlagEnv:                *FlagEnv,
			FlagCloudID:            *FlagCloudID,
			FlagAESKey:             *FlagAESKey,
			FlagRescheduleInterval: *FlagRescheduleInterval,
			FlagHeartbeatInterval:  *FlagHeartbeatInterval,
			Clients:                map[string]ClientConfig{},
		}
	}
	lock.Unlock()
	return conf
}
