package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

var conf *config

type collectionNames struct {
	CollectionCloud string
	CollectionUser  string
	CollectionFile  string
	MigrationAdvice string
	CollectionAK    string
}

type config struct {
	FlagMongo              string
	FlagAddress            string
	FlagEnv                string
	FlagCloudID            string
	FlagAESKey             string
	FlagRescheduleInterval time.Duration
	FlagHeartbeatInterval  time.Duration
	CollectionNames        collectionNames
}

type Config struct {
	*config
}

func SetGlobalConfig(mongo, address, env, cloudID, aesKey string, rescheduleInterval, heartbeatInterval time.Duration) {
	conf = &config{
		FlagMongo:              mongo,
		FlagAddress:            address,
		FlagEnv:                env,
		FlagCloudID:            cloudID,
		FlagAESKey:             aesKey,
		FlagRescheduleInterval: rescheduleInterval,
		FlagHeartbeatInterval:  heartbeatInterval,
		CollectionNames: collectionNames{
			CollectionCloud: "Cloud",
			CollectionUser:  "User",
			CollectionFile:  "File",
			MigrationAdvice: "MigrationAdvice",
			CollectionAK:    "AccessKey",
		},
	}
}

func GetConfig() Config {
	if conf == nil {
		logrus.Errorf("读取配置文件时，config 尚未设置！！")
	}
	return Config{conf}
}
