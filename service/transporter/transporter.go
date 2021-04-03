package main

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"strconv"
)

func main() {
	StartServe()
}

// 开启 transporter 服务
func StartServe() {
	flags := []cli.Flag{
		&cli.PathFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "JcsPan Transporter config file",
			Value:   "./transporter_config.json",
		},
	}
	app := cli.App{
		Name:    "Jcs-Transporter",
		Usage:   "Transporter backend for JcsPan",
		Version: util.GetVersionStr(),
		Authors: []*cli.Author{&cli.Author{
			Name:  "Zhang Junhua",
			Email: "zhangjh@mail.act.buaa.edu.cn",
		}},
		Flags: flags,
		Action: func(c *cli.Context) error {
			configFilePath := c.Path("config")
			err := util.ReadConfigFromFile(configFilePath)
			if err != nil {
				logrus.Errorf("Read config file fail:%v", err)
				return err
			}
			router, _ := initRouterAndProcessor()
			logrus.Infof("Transporter Started v%v at: %v:%v", util.GetVersionStr(), util.CONFIG.Host, util.CONFIG.Port)
			logrus.Info(http.ListenAndServe(":"+strconv.Itoa(util.CONFIG.Port), router))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func initRouterAndProcessor() (*controller.Router, *controller.TaskProcessor) {
	var storage model.TaskStorage
	var clientDatabase model.StorageDatabase
	var fileDatabase model.FileDatabase
	if util.CONFIG.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if util.CONFIG.Database.Driver == util.MongoDB {
		storage, _ = model.NewMongoTaskStorage()
		clientDatabase, _ = model.NewMongoStorageDatabase()
		fileDatabase, _ = model.NewMongoFileDatabase()
	} else {
		storage = model.NewInMemoryTaskStorage()
		clientDatabase = model.NewSimpleInMemoryStorageDatabase()
		fileDatabase = model.NewInMemoryFileDatabase()
	}
	processor := controller.TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(clientDatabase)
	// 初始化 FileInfo 数据库
	processor.FileDatabase = fileDatabase
	// 初始化路由
	router := controller.NewRouter(processor)
	// 启动 processor
	processor.StartProcessTasks(context.Background())

	return router, &processor
}
