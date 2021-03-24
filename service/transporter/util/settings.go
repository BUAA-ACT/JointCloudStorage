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
	DebugMode:            true,
	Database:             InMemoryDB,
	UploadFileTempPath:   "./tmp/upload/",
	DownloadFileTempPath: "./tmp/download/",
}

type Configuration struct {
	DebugMode            bool
	Database             string
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
	if CONFIG.Database != InMemoryDB && CONFIG.Database != MongoDB {
		err = errors.New("nonsupport database type")
	}
	return
}
