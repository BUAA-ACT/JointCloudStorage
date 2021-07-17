package main

import (
	"flag"
	"fmt"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/keySyn"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/newcloud"
)

const (
	Version = "v0.2"
)

var (
	flagMongo              = flag.String("mongo", "mongodb://localhost:27017", "mongodb address")
	flagAddress            = flag.String("addr", ":8082", "scheduler address")
	flagEnv                = flag.String("env", "test", "dev|test|prod")
	flagCloudID            = flag.String("cid", "aliyun-beijing", "cloud id")
	flagAESKey             = flag.String("aes", "1234567890123456", "aes key")
	flagRescheduleInterval = flag.Duration("reschedule", time.Minute*1, "reschedule interval")
	flagHeartbeatInterval  = flag.Duration("heartbeat", time.Second*30, "heartbeat interval")

	db      dao.Database
	addrMap = make(map[string]string)
)

func FlagParse(env string) {
	if env == "debug" {
		flagMongo = flag.String("mongo", "mongodb://192.168.105.8:20100", "mongodb address")
		flagEnv = flag.String("env", "dev", "Database name used for Clouds storage.")
	}
	flag.Parse()
}

func Init() {

	// Set logging format
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	// 初始化全局设置
	config.SetGlobalConfig(*flagMongo, *flagAddress, *flagEnv, *flagCloudID, *flagAESKey, *flagRescheduleInterval, *flagHeartbeatInterval)

	// Init DAO instance
	var err error
	db = dao.GetDatabaseInstance()
	if err != nil {
		panic(err)
	}

	// Init address map
	clouds, err := db.GetAllClouds()
	if err != nil {
		panic(err)
	}
	for _, c := range clouds {
		addrMap[c.CloudID] = c.Address
	}

	// Switch to release mode
	// if *flagEnv == "prod" {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
}

func NewRouter(r *gin.Engine) {
	r.GET("/storage_plan", GetStoragePlan)
	r.GET("/download_plan", GetDownloadPlan)
	r.GET("/status", GetStatus)
	r.GET("/all_clouds_status", GetAllCloudsStatus)

	r.POST("/storage_plan", PostStoragePlan)
	r.POST("/metadata", PostMetadata)
	r.POST("/update_clouds", PostUpdateClouds)
}

func main() {
	fmt.Println("this is main func")
	FlagParse("")
	Init()
	log.Infoln("Starting scheduler", Version)

	r := gin.Default()
	NewRouter(r)
	newcloud.Router(r, *flagMongo, *flagEnv, *flagCloudID, "production")
	keySyn.KeySynInit(*flagCloudID, r)
	go reSchedule(*flagRescheduleInterval)
	go heartbeat(*flagHeartbeatInterval)

	r.Run(*flagAddress)
}
