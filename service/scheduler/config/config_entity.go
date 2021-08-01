package config

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type collectionNames struct {
	CollectionCloud           string
	CollectionUser            string
	CollectionFile            string
	CollectionMigrationAdvice string
	CollectionAccessKey       string
	CollectionTempCloud       string
	CollectionVoteCloud       string
}

type Config struct {
	FlagMongo              string
	FlagAddress            string
	FlagEnv                string
	FlagCloudID            string
	FlagAESKey             string
	FlagRescheduleInterval time.Duration
	FlagHeartbeatInterval  time.Duration
}

func Factory() *Config {
	return &Config{
		FlagMongo:              "",
		FlagAddress:            "",
		FlagEnv:                "",
		FlagCloudID:            "",
		FlagAESKey:             "",
		FlagRescheduleInterval: 0,
		FlagHeartbeatInterval:  0,
	}
}
