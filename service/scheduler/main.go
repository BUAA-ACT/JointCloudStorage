package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/dao"
)

const (
	Version = "v0.2"

	CollectionCloud = "Cloud"
	CollectionUser  = "User"
	CollectionFile  = "File"
	MigrationAdvice = "MigrationAdvice"
)

var (
	flagMongo              = flag.String("mongo", "mongodb://localhost:27017", "mongodb address")
	flagAddress            = flag.String("addr", ":8082", "scheduler address")
	flagEnv                = flag.String("env", "test", "dev|test|prod")
	flagCloudID            = flag.String("cid", "aliyun-beijing", "cloud id")
	flagAESKey             = flag.String("aes", "1234567890123456", "aes key")
	flagRescheduleInterval = flag.Duration("reschedule", time.Hour*24*30, "reschedule interval")
	flagHeartbeatInterval  = flag.Duration("heartbeat", time.Second*10, "heartbeat interval")

	db      *dao.Dao
	addrMap = make(map[string]string)
)

func init() {
	flag.Parse()

	// Set logging format
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	// Init DAO instance
	var err error
	db, err = dao.NewDao(*flagMongo, *flagEnv, CollectionCloud, CollectionUser, CollectionFile, MigrationAdvice)
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

func main() {
	log.Infoln("Starting scheduler", Version)

	r := gin.Default()

	r.GET("/storage_plan", GetStoragePlan)
	r.GET("/download_plan", GetDownloadPlan)
	r.GET("/status", GetStatus)

	r.POST("/storage_plan", PostStoragePlan)
	r.POST("/metadata", PostMetadata)

	go reSchedule(*flagRescheduleInterval)
	go heartbeat(*flagHeartbeatInterval)

	r.Run(*flagAddress)
}
