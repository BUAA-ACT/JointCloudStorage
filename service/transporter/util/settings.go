package util

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	InMemoryDB  = "InMemoryDB"
	MongoDB     = "MongoDB"
	MinioClient = "Minio"
	AwsS3Client = "Aws"
)

var Config = Configuration{
	DebugMode: true,
	Database: DBConfiguration{
		Driver:       MongoDB,
		Host:         "192.168.105.8",
		Port:         "20100",
		DatabaseName: "dev",
	},
	UploadFileTempPath:   "./tmp/upload/",
	DownloadFileTempPath: "./tmp/download/",
	DefaultStorageClient: AwsS3Client,
	Port:                 8083,
	Host:                 "0.0.0.0",
	ZookeeperHost:        "192.168.105.13",
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
	DefaultStorageClient string
	Port                 int
	Host                 string
	ZookeeperHost        string
	LocalCloudID         string
	SchedulerHost        string
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
	Config = conf
	err = CheckConfig()
	return err
}

func CheckConfig() (err error) {
	if IsDir(Config.UploadFileTempPath) {
		err = os.MkdirAll(Config.UploadFileTempPath, os.ModePerm)
	}
	if IsDir(Config.DownloadFileTempPath) {
		err = os.MkdirAll(Config.DownloadFileTempPath, os.ModePerm)
	}
	if Config.Database.Driver != InMemoryDB && Config.Database.Driver != MongoDB {
		err = errors.New("nonsupport database type")
	}
	return
}
