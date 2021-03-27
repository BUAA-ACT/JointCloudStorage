package util

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	InMemoryDB = "InMemoryDB"
	MongoDB    = "MongoDB"
)

var CONFIG = Configuration{
	DebugMode: true,
	Database: DBConfiguration{
		Driver:       MongoDB,
		Host:         "192.168.105.8",
		Port:         "20100",
		DatabaseName: "dev",
	},
	UploadFileTempPath:   "./tmp/upload/",
	DownloadFileTempPath: "./tmp/download/",
}

type DBConfiguration struct {
	Driver       string
	Host         string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}

type Configuration struct {
	DebugMode            bool
	Database             DBConfiguration
	UploadFileTempPath   string
	DownloadFileTempPath string
}

func ReadConfigFromFile(configFilepath string) error {
	file, err := os.Open(configFilepath)
	if err != nil {
		return errors.New("open config file fail")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	}
	CONFIG = conf
	err = CheckConfig()
	return err
}

func CheckConfig() (err error) {
	if IsDir(CONFIG.UploadFileTempPath) {
		err = os.MkdirAll(CONFIG.UploadFileTempPath, os.ModePerm)
	}
	if IsDir(CONFIG.DownloadFileTempPath) {
		err = os.MkdirAll(CONFIG.DownloadFileTempPath, os.ModePerm)
	}
	if CONFIG.Database.Driver != InMemoryDB && CONFIG.Database.Driver != MongoDB {
		err = errors.New("nonsupport database type")
	}
	return
}
