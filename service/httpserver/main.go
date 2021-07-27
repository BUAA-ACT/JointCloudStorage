package main

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/controller"
	"cloud-storage-httpserver/dao"
	"cloud-storage-httpserver/service/sweeper"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"strconv"
)

func setupRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		f := c.Query("test")
		fmt.Println(f)
		c.String(http.StatusOK, "pong")
	})

	/* test function */
	r.POST("/testPost", controller.UserTestPost)
	r.GET("/testGet", controller.UserTestGet)
	r.GET("/cookie", controller.CookieTestGet)
	r.POST("/header", controller.HeaderTestPost)
	r.GET("/websocket", controller.WebSocketTestGet)

	/* user function */
	r.POST("/user/register", controller.UserRegister)
	r.POST("/user/checkVerifyCode", controller.UserCheckVerifyCode)
	r.POST("/user/login", controller.UserLogin)
	r.POST("/user/logout", controller.UserLogout)
	r.POST("/user/checkValidity", controller.UserCheckValidity)
	r.POST("/user/changePassword", controller.UserChangePassword)
	r.POST("/user/changeEmail", controller.UserChangeEmail)
	r.POST("/user/changeNickname", controller.UserChangeNickname)
	r.POST("/user/getUserInfo", controller.UserGetInfo)
	r.POST("/user/changeUserPreference", controller.UserSetPreference)
	r.POST("/user/newKey", controller.UserAddKey)
	r.POST("/user/getUserKeys", controller.UserGetKeys)
	r.POST("/user/deleteKey", controller.UserDeleteKey)
	r.POST("/user/changeKeyStatus", controller.UserChangeKeyStatus)
	r.POST("/user/changeKeyComment", controller.UserChangeKeyComment)
	r.POST("/user/remakeKey", controller.UserRemakeKey)
	//r.POST("/user/uploadAvatar",controller.UserUploadAvatar)

	/* cloud function */
	r.POST("/cloud/thisName", controller.GetThisCloudName)

	/* plan function */
	//r.POST("/plan/changeUserStoragePlan", controller.UserSetStoragePlan)
	r.POST("/plan/chooseStoragePlan", controller.UserChooseStoragePlan)
	r.POST("/plan/acceptAdvice", controller.UserAcceptStoragePlan)
	r.POST("/plan/getAllStoragePlan", controller.UserGetAllStoragePlan)
	r.POST("/plan/getNewAdvice", controller.UserGetAdvice)
	r.POST("/plan/abandonAdvice", controller.UserAbandonAdvice)

	/* file function */
	r.POST("/file/getFiles", controller.UserGetFiles)
	r.POST("/file/changeFilePath", controller.UserChangeFilePath)
	r.POST("/file/changeFileName", controller.UserChangeFileName)
	r.POST("/file/preUploadFile", controller.UserPreUploadFile)
	r.POST("/file/downloadFile", controller.UserDownloadFile)
	r.POST("/file/deleteFile", controller.UserDeleteFile)

	/* task function */
	r.POST("/task/getTask", controller.UserGetTask)
	r.GET("/task/getMigration", controller.UserGetMigrationTask)

	/* admin & cloud function */
	r.POST("/cloud/getAllClouds", controller.GetAllClouds)
	r.POST("/cloud/newCloud", controller.AdminAddCloud)
	r.POST("/cloud/changeCloudInfo", controller.AdminChangeCloudInfo)
	r.POST("/cloud/vote", controller.AdminVoteForCloud)
	r.POST("/cloud/getVoteRequests", controller.AdminGetVoteRequests)
	r.POST("/cloud/getAddedClouds", controller.AdminGetAddedClouds)
}

func StartServe(configFilePath string) {
	args.LoadProperties(configFilePath)
	r := gin.Default()
	//cross
	r.Use(controller.Cors)
	//route
	setupRouter(r)
	dao.ConnectInitDao()
	dao.AddIndex()
	go sweeper.CleanAccessToken()
	go sweeper.CleanVerifyCode()
	_ = r.Run(":" + strconv.FormatUint(*args.HttpserverPort, 10))
}

func main() {
	flags := []cli.Flag{
		&cli.PathFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "JcsPan httpserver config file",
			Value:   "./httpserver.properties",
		},
	}
	app := cli.App{
		Name:  "Jcs-Httpserver",
		Usage: "Httpserver backend for JcsPan",
		Authors: []*cli.Author{&cli.Author{
			Name:  "Zhang Junhua",
			Email: "zhangjh@mail.act.buaa.edu.cn",
		}},
		Flags: flags,
		Action: func(c *cli.Context) error {
			configFilePath := c.Path("config")
			StartServe(configFilePath)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		print(err)
	}
}
