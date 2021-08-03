package main

import (
	"flag"
	"fmt"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/dao"
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
	config.GetConfig()

	//Switch to release mode
	//if *flagEnv == "prod" {
	//	gin.SetMode(gin.ReleaseMode)
	//}

}

func serverPlugIn(r *gin.Engine) {
	server.IDInit(*config.FlagCloudID)
	// server module use the databases below
	databaseMap := map[string]map[string]*dao.CollectionConfig{
		*config.FlagEnv: {
			config.CloudCollectionName:           nil,
			config.UserCollectionName:            nil,
			config.FileCollectionName:            nil,
			config.MigrationAdviceCollectionName: nil,
		},
	}

	err := server.DaoInit(*config.FlagMongo, databaseMap)
	if err != nil {
		log.Errorf("server plug in failed with error : %s", err.Error())
	}
	server.SetCloudCol(databaseMap[*config.FlagEnv][config.CloudCollectionName].CollectionHandler)
	server.SetUserCol(databaseMap[*config.FlagEnv][config.UserCollectionName].CollectionHandler)
	server.SetFileCol(databaseMap[*config.FlagEnv][config.FileCollectionName].CollectionHandler)
	server.SetAdviceCol(databaseMap[*config.FlagEnv][config.MigrationAdviceCollectionName].CollectionHandler)
	server.RouteInit(r)
}

func newCloudPlugIn(r *gin.Engine) {
	newcloud.IDInit(*config.FlagCloudID, *config.FlagEnv)
	// new cloud module use the databases below
	databaseMap := map[string]map[string]*dao.CollectionConfig{
		*config.FlagEnv: {
			config.CloudCollectionName:     nil,
			config.TempCloudCollectionName: nil,
			config.VoteCloudCollectionName: nil,
		},
	}
	err := newcloud.DaoInit(*config.FlagMongo, databaseMap)
	if err != nil {
		log.Errorf("server plug in failed with error : %s", err.Error())
	}
	newcloud.SetCloudCol(databaseMap[*config.FlagEnv][config.CloudCollectionName].CollectionHandler)
	newcloud.SetTempCloudCol(databaseMap[*config.FlagEnv][config.TempCloudCollectionName].CollectionHandler)
	newcloud.SetVoteCloudCol(databaseMap[*config.FlagEnv][config.VoteCloudCollectionName].CollectionHandler)
	newcloud.RouteInit(r)
}

func keySynPlugIn(r *gin.Engine) {
	keySyn.IDInit(*config.FlagCloudID)
	// key synchronize module use the databases below
	databaseMap := map[string]map[string]*dao.CollectionConfig{
		*config.FlagEnv: map[string]*dao.CollectionConfig{
			config.AccessKeyCollectionName: nil,
			config.CloudCollectionName:     nil,
		},
	}
	err := keySyn.DaoInit(*config.FlagMongo, databaseMap)
	if err != nil {
		log.Errorf("server plug in failed with error : %s", err.Error())
	}
	keySyn.SetKeyCol(databaseMap[*config.FlagEnv][config.AccessKeyCollectionName].CollectionHandler)
	keySyn.SetCloudCol(databaseMap[*config.FlagEnv][config.CloudCollectionName].CollectionHandler)
	keySyn.RouteInit(r)
}

func main() {
	fmt.Println("this is main func")
	FlagParse("")
	Init()
	log.Infoln("Starting scheduler", config.Version)

	r := gin.Default()
	serverPlugIn(r)
	keySynPlugIn(r)
	newCloudPlugIn(r)

	go server.ReSchedule(*config.FlagRescheduleInterval)
	go server.Heartbeat(*config.FlagHeartbeatInterval)

	r.Run(*config.FlagAddress)
}
