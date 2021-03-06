package main

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/storageInterface"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	g errgroup.Group
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
			router, _, jsi := initRouterAndProcessor()
			g.Go(func() error {
				logrus.Infof("Transporter Started v%v at: %v:%v", util.GetVersionStr(), util.Config.Host, util.Config.Port)
				return http.ListenAndServe(":"+strconv.Itoa(util.Config.Port), router)
			})

			g.Go(func() error {
				logrus.Infof("JSI Started v%v at: %v:%v", util.GetVersionStr(), util.Config.Host, util.Config.JSIPort)
				return http.ListenAndServe(":"+strconv.Itoa(util.Config.JSIPort), jsi)
			})

			if err := g.Wait(); err != nil {
				return err
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func initRouterAndProcessor() (*controller.Router, *controller.TaskProcessor, *storageInterface.JointStorageInterface) {
	var storage model.TaskStorage
	var clientDatabase model.CloudDatabase
	var fileDatabase model.FileDatabase
	if util.Config.DebugMode {
		logrus.Warning("Run in debug mode")
		switch util.Config.DebugLevel {
		case "Trace":
			logrus.SetLevel(logrus.TraceLevel)
			logrus.Warning("DebugLevel: Trace")
		default:
			logrus.SetLevel(logrus.DebugLevel)
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if util.Config.Database.Driver == util.MongoDB {
		storage, _ = model.NewMongoTaskStorage()
		clientDatabase, _ = model.NewMongoCloudDatabase()
		fileDatabase, _ = model.NewMongoFileDatabase()
	} else {
		storage = model.NewInMemoryTaskStorage()
		clientDatabase = model.NewSimpleInMemoryStorageDatabase()
		fileDatabase = model.NewInMemoryFileDatabase()
	}
	processor := controller.TaskProcessor{}
	// 初始化 Dao
	var err error
	processor.Dao, err = model.InitDao()
	if err != nil {
		util.Log(logrus.FatalLevel, "init Processor", "error init mongodb connection", "Dao", "err", err.Error())
	}

	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(clientDatabase)
	// 初始化 FileInfo 数据库
	processor.FileDatabase = fileDatabase
	// 初始化 lock
	lock, _ := controller.NewLock(util.Config.ZookeeperHost)
	processor.Lock = lock
	//processor.lock.UnLockAll("/tester")
	// 初始化 scheduler
	scheduler := controller.JcsPanScheduler{
		LocalCloudID:     util.Config.LocalCloudID,
		SchedulerHostUrl: util.Config.SchedulerHost,
		ReloadCloudInfo:  true,
		CloudDatabase:    clientDatabase,
	}
	processor.Scheduler = &scheduler
	// 初始化 Monitor
	userDB, _ := model.NewMongoUserDatabase()
	processor.Monitor = controller.NewTrafficMonitor(userDB)
	processor.UserDatabase = userDB
	// 初始化 tempFile
	tfs, _ := util.NewTempFileStorage(util.Config.DownloadFileTempPath, time.Hour*8)
	processor.TempFileStorage = tfs
	// 初始化 AccessKeyDB
	dao, _ := model.InitDao()
	accessKeyDB := model.AccessKeyDB{Dao: dao}
	processor.AccessKeyDatabase = &accessKeyDB

	// 初始化路由
	router := controller.NewRouter(processor)
	// 初始化 JSI
	jsi := storageInterface.NewInterface(&processor)
	// 启动 processor
	processor.StartProcessTasks(context.Background())

	return router, &processor, jsi
}
