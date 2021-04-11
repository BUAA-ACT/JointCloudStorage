package util

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
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
	LocalCloudID:         "aliyun-beijing",
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
	if Config.UploadFileTempPath[len(Config.UploadFileTempPath)-1] != '/' {
		Config.UploadFileTempPath += "/"
	}
	if Config.DownloadFileTempPath[len(Config.DownloadFileTempPath)-1] != '/' {
		Config.DownloadFileTempPath += "/"
	}

	if !IsDir(Config.UploadFileTempPath) {
		Log(logrus.InfoLevel, "CheckConfig", "Create upload Dir", "", "", Config.UploadFileTempPath)
		err = os.MkdirAll(Config.UploadFileTempPath, os.ModePerm)
	}
	if !IsDir(Config.DownloadFileTempPath) {
		Log(logrus.InfoLevel, "CheckConfig", "Create download Dir", "", "", Config.UploadFileTempPath)
		err = os.MkdirAll(Config.DownloadFileTempPath, os.ModePerm)
	}
	if Config.Database.Driver != InMemoryDB && Config.Database.Driver != MongoDB {
		err = errors.New("nonsupport database type")
	}
	return
}
