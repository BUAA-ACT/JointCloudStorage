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

type CollectionConfig struct {
	CollectionHandler *mongo.Collection
}

type DatabaseConfig struct {
	DatabaseHandler *mongo.Database
	Collections     map[string]CollectionConfig
}

type ClientConfig struct {
	Client    *mongo.Client
	Databases map[string]DatabaseConfig
}

type Config struct {
	FlagMongo              string
	FlagAddress            string
	FlagEnv                string
	FlagCloudID            string
	FlagAESKey             string
	FlagRescheduleInterval time.Duration
	FlagHeartbeatInterval  time.Duration
	Clients                map[string]ClientConfig
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
