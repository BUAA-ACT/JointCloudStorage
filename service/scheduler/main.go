package main

import (
	"flag"
	"fmt"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/keySyn"
	"shaoliyin.me/jcspan/newcloud"
	"shaoliyin.me/jcspan/server"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func FlagParse(env string) {
	if env == "debug" {
		config.FlagMongo = flag.String("mongo", "mongodb://192.168.105.8:20100", "mongodb address")
		config.FlagEnv = flag.String("env", "dev", "Database name used for Clouds storage.")
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
	config.SetGlobalConfig()
	//
	//// Init DAO instance
	//var err error
	//db = dao.GetDatabaseInstance()
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Init address map
	//clouds, err := db.GetAllClouds()
	//if err != nil {
	//	panic(err)
	//}
	//for _, c := range clouds {
	//	addrMap[c.CloudID] = c.Address
	//}
	//
	//// Switch to release mode
	//// if *flagEnv == "prod" {
	//// 	gin.SetMode(gin.ReleaseMode)
	//// }
	//
	//
	//
	//localMongo, err := dao.NewDao(mongo, databasename, CollectionCloud, CollectionUser, CollectionFile, MigrationAdvice, "")
	//if err != nil {
	//	return err
	//}
	//
	//localMongoTempCloud, err = dao.NewDao(mongo, databasename, CollectionTempCloud, CollectionUser, CollectionFile, MigrationAdvice, "")
	//if err != nil {
	//	return err
	//}
	//
	//localMongoVoteRequest, err = dao.NewDao(mongo, databasename, CollectionVoteCloud, CollectionUser, CollectionFile, MigrationAdvice, "")
	//if err != nil {
	//	return err
	//}
	//
	//localid = cid
	//tempNotFound = errors.New("TempCloud not Found.")
}

func main() {
	fmt.Println("this is main func")
	FlagParse("")
	Init()
	log.Infoln("Starting scheduler", config.Version)

	r := gin.Default()
	server.NewRouter(r)
	newcloud.PlugIn(r, *flagMongo, *flagEnv, *flagCloudID, "production")
	keySyn.KeySyncInit(*flagCloudID, r)
	go server.ReSchedule(*config.FlagRescheduleInterval)
	go server.Heartbeat(*config.FlagHeartbeatInterval)
	//
	r.Run(*config.FlagAddress)
}
