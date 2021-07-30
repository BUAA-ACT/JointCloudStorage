package server

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"shaoliyin.me/jcspan"
	"shaoliyin.me/jcspan/newcloud"
	"shaoliyin.me/jcspan/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shaoliyin.me/jcspan/dao"
)

const (
	codeOK            = 200
	codeBadRequest    = 400
	codeUnauthorized  = 401
	codeInternalError = 500

	ReplicaMode = "Replica"
	ECMode      = "EC"
)

var (
	errorMsg = map[int]string{
		codeOK:            "OK",
		codeBadRequest:    "Bad Request",
		codeUnauthorized:  "Unauthorized",
		codeInternalError: "Internal Server Error",
	}
)

func NewRouter(r *gin.Engine) {
	r.GET("/storage_plan", GetStoragePlan)
	r.GET("/download_plan", GetDownloadPlan)
	r.GET("/status", GetStatus)
	r.GET("/all_clouds_status", GetAllCloudsStatus)

	r.POST("/storage_plan", PostStoragePlan)
	r.POST("/metadata", PostMetadata)
	r.POST("/update_clouds", PostUpdateClouds)
}

func GetStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param GetStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		main.logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	main.logInfo("Receive GetStoragePlan", requestID, param)

	clouds, err := db.GetAllClouds()
	if err != nil {
		main.logError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// 计算最佳方案
	plan := main.storagePlan(param, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      plan,
	})

	main.logInfo("Response GetStoragePlan", requestID, plan)
}

func GetDownloadPlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param GetDownloadPlanParam
	err := c.BindJSON(&param)
	if err != nil {
		main.logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	main.logInfo("Receive GetDownloadPlan", requestID, param)

	user, err := db.GetUser(param.UserID)
	if err != nil {
		main.logError(err, requestID, "GetUserInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	clouds, err := db.GetAllClouds()
	if err != nil {
		main.logError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// 计算最佳方案
	plan := main.downloadPlan(user.StoragePlan, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      user.StoragePlan,
	})

	main.logInfo("Response GetDownloadPlan", requestID, plan)
}

func GetStatus(c *gin.Context) {
	requestID := uuid.New().String()
	var param GetStatusParam
	err := c.BindJSON(&param)
	if err != nil {
		main.logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	main.logTrace("Receive GetStatus", requestID, param)

	// 验证请求来源是否合法
	_, err = db.GetCloud(param.CloudID)
	if err != nil {
		main.logError(err, requestID, "GetCloudInfo failed", param.CloudID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"RequestID": requestID,
			"Code":      codeUnauthorized,
			"Msg":       errorMsg[codeUnauthorized],
		})
		return
	}

	// 获取本云信息
	cloud, err := db.GetCloud(*flagCloudID)
	if err != nil {
		main.logError(err, requestID, "GetCloudInfo failed", *flagCloudID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      cloud,
	})

	main.logTrace("Response GetStatus", requestID, cloud)
}

func PostStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param PostStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		main.logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	main.logInfo("Receive PostStoragePlan", requestID, param)

	if param.CloudID == *flagCloudID {
		// 来自本云httpserver的请求
		plan := &param.StoragePlan
		plan.StoragePrice = main.calStoragePrice(*plan)
		plan.Availability = main.calAvailability(*plan)
		plan.TrafficPrice = main.calTrafficPrice(*plan, false)
		var users []dao.AccessCredential
		ch := make(chan *dao.AccessCredential)

		clouds, err := db.GetAllClouds()
		if err != nil {
			main.logError(err, requestID, "GetAllCloudInfo failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, cloud := range clouds {
			// 通知其他可用云
			if cloud.CloudID == param.CloudID {
				continue
			}
			go func(cloud dao.Cloud) {
				u, err := main.sendPostStoragePlan(param, cloud.CloudID)
				if err != nil {
					main.logError(err, requestID, "sendPostStoragePlan failed", param, cloud)
				}
				ch <- u
			}(cloud)
		}

		for i := 0; i < len(clouds)-1; i++ {
			u := <-ch
			if u != nil {
				users = append(users, *u)
			}
		}

		if len(users) < len(clouds)-1 {
			main.logError(nil, requestID, "Some sendPostStoragePlan failed", len(clouds)-1, len(users))
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
		}

		// 更正成本计算
		var user dao.User
		user, err = db.GetUser(param.UserID)
		if err != nil {
			main.logError(err, requestID, "更新存储方案时获取用户信息失败")
		} else {
			user.StoragePlan = param.StoragePlan
			user.Preference.Vendor = param.StoragePlan.N // 存储偏好，副本数
			user.Preference.StoragePrice = math.Max(user.Preference.StoragePrice, user.StoragePlan.StoragePrice)
			user.Preference.TrafficPrice = math.Max(user.Preference.TrafficPrice, user.StoragePlan.TrafficPrice)
			user.Preference.Availability = math.Min(user.Preference.Availability, user.StoragePlan.Availability)
			err = db.InsertUser(user)
			if err != nil {
				main.logError(err, requestID, "更新用户存储方案失败")
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
			"Data":      users,
		})
		main.logInfo("Response PostStoragePlan", requestID, users)
	} else {
		// 来自其他云scheduler的请求
		_, err = db.GetCloud(param.CloudID)
		if err != nil {
			main.logError(err, requestID, "GetCloudInfo failed", param.CloudID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"RequestID": requestID,
				"Code":      codeUnauthorized,
				"Msg":       errorMsg[codeUnauthorized],
			})
			return
		}

		// 新建用户
		passwd := param.Password
		// 密码为空时从数据库获取已有 User
		var user dao.User
		if passwd == "" {
			user, err = db.GetUser(param.UserID)
			if err != nil {
				main.logError(err, requestID, "Get user failed", user)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
			user.StoragePlan = param.StoragePlan
			user.LastModified = time.Now()
		} else {
			user = dao.User{
				UserId:       param.UserID,
				Email:        param.UserID,
				Nickname:     param.UserID,
				Password:     tools.AesEncrypt(passwd, *flagAESKey),
				Role:         dao.RoleGuest,
				LastModified: time.Now(),
				StoragePlan:  param.StoragePlan,
			}
		}
		err = db.InsertUser(user)
		if err != nil {
			main.logError(err, requestID, "InsertUser failed", user)
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		cred := dao.AccessCredential{
			CloudID:  *flagCloudID,
			UserID:   param.UserID,
			Password: passwd,
		}
		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
			"Data":      []dao.AccessCredential{cred},
		})
		main.logInfo("Response PostStoragePlan", requestID, []dao.User{user})
	}
}

func PostMetadata(c *gin.Context) {
	requestID := uuid.New().String()

	var param PostMetadataParam
	err := c.BindJSON(&param)
	if err != nil {
		main.logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	main.logInfo("Receive PostMetadata", requestID, param)

	if param.CloudID == *flagCloudID {
		// 来自本云httpserver的请求
		var errs []error
		var routine int
		ch := make(chan error)

		clouds, err := db.GetAllClouds()
		if err != nil {
			main.logError(err, requestID, "GetAllCloudInfo failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, cloud := range clouds {
			// 通知其他云
			if cloud.CloudID == param.CloudID {
				continue
			}
			routine++
			go func(cloud dao.Cloud) {
				err := main.sendPostMetadata(param, cloud.CloudID)
				if err != nil {
					main.logError(err, requestID, "sendPostMetadata failed", param, cloud)
				}
				ch <- err
			}(cloud)
		}

		for i := 0; i < routine; i++ {
			err := <-ch
			if err == nil {
				errs = append(errs, err)
			}
		}

		if len(errs) < routine {
			main.logError(nil, requestID, "Some sendPostMetadata failed", routine, len(errs))
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
		})
		main.logInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	} else {
		// 来自其他云scheduler的请求

		// 校验请求合法性
		_, err = db.GetCloud(param.CloudID)
		if err != nil {
			main.logError(err, requestID, "GetCloudInfo failed", param.CloudID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"RequestID": requestID,
				"Code":      codeUnauthorized,
				"Msg":       errorMsg[codeUnauthorized],
			})
			return
		}

		if param.Type == "Upload" {
			// 写入文件元信息
			err = db.InsertFiles(param.Files)
			if err != nil {
				main.logError(err, requestID, "InsertFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// 修改用户存储量
			err = db.ChangeVolume(param.UserID, "Upload", param.Files)
			if err != nil {
				main.logError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
		} else if param.Type == "Delete" {
			// 删除文件元信息
			err = db.DeleteFiles(param.Files)
			if err != nil {
				main.logError(err, requestID, "DeleteFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// 修改用户存储量
			err = db.ChangeVolume(param.UserID, "Delete", param.Files)
			if err != nil {
				main.logError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
		} else if param.Type == "Migrate" {
			// 元数据已在重新PostStoragePlan时修改，此时不需要做任何事
			// err = db.DeleteUser(param.UserID)
			// if err != nil {
			// 	logError(err, requestID, "DeleteUser failed", param.UserID)
			// 	c.JSON(http.StatusInternalServerError, gin.H{
			// 		"RequestID": requestID,
			// 		"Code":      codeInternalError,
			// 		"Msg":       errorMsg[codeInternalError],
			// 	})
			// 	return
			// }
		} else {
			main.logError(err, requestID, errorMsg[codeBadRequest])
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeBadRequest,
				"Msg":       errorMsg[codeBadRequest],
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
		})
		main.logInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	}
}

func Heartbeat(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		<-t.C
		requestID := uuid.New().String()
		main.logTrace("Starting to send heartbeat packages", requestID)

		clouds, err := db.GetOtherClouds(*flagCloudID)
		if err != nil {
			main.logError(err, requestID, "GetOtherClouds failed", *flagCloudID)
			continue
		}

		ch := make(chan error)
		param := GetStatusParam{CloudID: *flagCloudID}
		for _, cloud := range clouds {
			go func(cloud dao.Cloud) {
				c, err := main.sendGetStatus(param, cloud.CloudID)
				if err != nil {
					main.logError(err, requestID, "sendGetStatus failed", param, cloud)
					ch <- err
					return
				}
				err = db.UpdateCloud(*c)
				if err != nil {
					main.logError(err, requestID, "UpdateCloud failed", *c)
					ch <- err
					return
				}
				ch <- nil
			}(cloud)
		}

		var success int
		for range clouds {
			err := <-ch
			if err == nil {
				success++
			}
		}

		main.logTrace("Heartbeat finished", requestID, len(clouds), success)
	}
}

func GetAllCloudsStatus(c *gin.Context) {
	requestID := uuid.New().String()

	//查询所有的clouds
	clouds, err := db.GetAllClouds()
	main.logInfo("clouds:", requestID, clouds)
	if err != nil {
		//查询出错，报告错误
		main.logError(err, requestID, "query for all clouds status failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}
	//查询成功，返回数据
	//隐藏accesskey和secretkey
	for index, _ := range clouds {
		clouds[index].AccessKey = ""
		clouds[index].SecretKey = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      clouds,
	})
	return

}

func PostUpdateClouds(c *gin.Context) {
	requestID := uuid.New().String()
	//get the clouds
	main.logInfo("UpdateClouds:", requestID, c)
	var cloud dao.Cloud
	if err := c.ShouldBindJSON(&cloud); err != nil {
		//can't get the clouds
		//return the error
		main.logError(err, requestID, "can't get the clouds from paramators")
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}

	//update the clouds
	if err := db.UpdateCloud(cloud); err != nil {
		//log the error
		main.logError(err, requestID, "can't update the cloud")
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	//向其他云同步
	if c.GetHeader("Caller") == "http-server" {
		clouds, err := db.GetAllClouds()
		if err != nil {
			main.logError(err, requestID, "can't get other clouds")
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		b, err := json.Marshal(cloud)
		if err != nil {
			main.logError(err, requestID, "can't Marshal the cloud")
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, otherCLoud := range clouds {
			if otherCLoud.CloudID != *flagCloudID {
				body := bytes.NewBuffer(b)
				addr := main.genAddress(otherCLoud.CloudID, "/update_clouds")
				resp, err := http.Post(addr, "application/json", body)
				if err != nil || resp.StatusCode != 200 {
					main.logError(err, requestID, "can't syn to cloud: ", otherCLoud.CloudID)
					c.JSON(http.StatusBadRequest, gin.H{
						"RequestID": requestID,
						"Code":      codeInternalError,
						"Msg":       errorMsg[codeInternalError],
					})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
	main.logInfo("update the clouds succeeded!", requestID, cloud)
}
