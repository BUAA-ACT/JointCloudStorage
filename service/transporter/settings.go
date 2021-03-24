package main

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	InMemoryDB = "InMemoryDB"
	MongoDB    = "MongoDB"
)

var CONFIG Configuration

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
	return nil
}
