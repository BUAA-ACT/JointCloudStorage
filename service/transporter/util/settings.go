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
	DebugMode:  true,
	DebugLevel: "Trace",
	Database: DBConfiguration{
		Driver:       MongoDB,
		Host:         "192.168.105.13",
		Port:         "27017",
		DatabaseName: "qingdao",
	},
	TempFilePath:         "./tmp/",
	UploadFileTempPath:   "./tmp/upload/",
	DownloadFileTempPath: "./tmp/download/",
	DefaultStorageClient: AwsS3Client,
	JSIPort:              8085,
	Port:                 8083,
	Host:                 "0.0.0.0",
	ZookeeperHost:        "192.168.105.13",
	SchedulerHost:        "http://192.168.105.13:8282",
	//SchedulerHost:        "http://127.0.0.1:8082",
	LocalCloudID: "aliyun-qingdao",
	EnableHttps:  true,
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
	DebugLevel           string
	Database             DBConfiguration
	TempFilePath         string
	UploadFileTempPath   string
	DownloadFileTempPath string
	DefaultStorageClient string
	JSIPort              int
	Port                 int
	Host                 string
	ZookeeperHost        string
	LocalCloudID         string
	SchedulerHost        string
	EnableHttps          bool
}

func ReadConfigFromFile(configFilepath string) error {
	file, err := os.Open(configFilepath)
	if err != nil {
		return errors.New("open config file fail")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config
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
