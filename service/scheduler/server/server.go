package server

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/entity"
	"shaoliyin.me/jcspan/tools"
	"shaoliyin.me/jcspan/utils"
	"time"
)

func GetStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param entity.GetStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		tools.LogError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	tools.LogInfo("Receive GetStoragePlan", requestID, param)

	clouds, err := dao.GetAllClouds(cloudCol)
	if err != nil {
		tools.LogError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// 计算最佳方案
	plan := storagePlan(param, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      plan,
	})

	tools.LogInfo("Response GetStoragePlan", requestID, plan)
}

func GetDownloadPlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param entity.GetDownloadPlanParam
	err := c.BindJSON(&param)
	if err != nil {
		tools.LogError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	tools.LogInfo("Receive GetDownloadPlan", requestID, param)

	user, err := dao.GetUser(userCol, param.UserID)
	if err != nil {
		tools.LogError(err, requestID, "GetUserInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	clouds, err := dao.GetAllClouds(cloudCol)
	if err != nil {
		tools.LogError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// 计算最佳方案
	plan := downloadPlan(user.StoragePlan, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      user.StoragePlan,
	})

	tools.LogInfo("Response GetDownloadPlan", requestID, plan)
}

func GetStatus(c *gin.Context) {
	requestID := uuid.New().String()
	var param entity.GetStatusParam
	err := c.BindJSON(&param)
	if err != nil {
		tools.LogError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	tools.LogTrace("Receive GetStatus", requestID, param)

	// 验证请求来源是否合法
	_, err = dao.GetCloud(cloudCol, param.CloudID)
	if err != nil {
		tools.LogError(err, requestID, "GetCloudInfo failed", param.CloudID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"RequestID": requestID,
			"Code":      codeUnauthorized,
			"Msg":       errorMsg[codeUnauthorized],
		})
		return
	}

	// 获取本云信息
	cloud, err := dao.GetCloud(cloudCol, *config.FlagCloudID)
	if err != nil {
		tools.LogError(err, requestID, "GetCloudInfo failed", *config.FlagCloudID)
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

	tools.LogTrace("Response GetStatus", requestID, cloud)
}

func PostStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param entity.PostStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		tools.LogError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	tools.LogInfo("Receive PostStoragePlan", requestID, param)

	if param.CloudID == *config.FlagCloudID {
		// 来自本云httpserver的请求
		plan := &param.StoragePlan
		plan.StoragePrice = calStoragePrice(*plan)
		plan.Availability = calAvailability(*plan)
		plan.TrafficPrice = calTrafficPrice(*plan, false)
		var users []entity.AccessCredential
		ch := make(chan *entity.AccessCredential)

		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			tools.LogError(err, requestID, "GetAllCloudInfo failed")
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
			go func(cloud entity.Cloud) {
				u, err := utils.SendPostStoragePlan(cloudCol, param, cloud.CloudID)
				if err != nil {
					tools.LogError(err, requestID, "sendPostStoragePlan failed", param, cloud)
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
			tools.LogError(nil, requestID, "Some sendPostStoragePlan failed", len(clouds)-1, len(users))
			//c.JSON(http.StatusInternalServerError, gin.H{
			//	"RequestID": requestID,
			//	"Code":      codeInternalError,
			//	"Msg":       errorMsg[codeInternalError],
			//})
		}

		// 更正成本计算
		var user entity.User
		user, err = dao.GetUser(userCol, param.UserID)
		if err != nil {
			tools.LogError(err, requestID, "更新存储方案时获取用户信息失败")
		} else {
			user.StoragePlan = param.StoragePlan
			user.Preference.Vendor = param.StoragePlan.N // 存储偏好，副本数
			user.Preference.StoragePrice = math.Max(user.Preference.StoragePrice, user.StoragePlan.StoragePrice)
			user.Preference.TrafficPrice = math.Max(user.Preference.TrafficPrice, user.StoragePlan.TrafficPrice)
			user.Preference.Availability = math.Min(user.Preference.Availability, user.StoragePlan.Availability)
			err = dao.InsertUser(userCol, user)
			if err != nil {
				tools.LogError(err, requestID, "更新用户存储方案失败")
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
			"Data":      users,
		})
		tools.LogInfo("Response PostStoragePlan", requestID, users)
	} else {
		// 来自其他云scheduler的请求
		_, err = dao.GetCloud(cloudCol, param.CloudID)
		if err != nil {
			tools.LogError(err, requestID, "GetCloudInfo failed", param.CloudID)
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
		var user entity.User
		if passwd == "" {
			user, err = dao.GetUser(userCol, param.UserID)
			if err != nil {
				tools.LogError(err, requestID, "Get user failed", user)
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
			user = entity.User{
				UserId:       param.UserID,
				Email:        param.UserID,
				Nickname:     param.UserID,
				Password:     tools.AesEncrypt(passwd, *config.FlagAESKey),
				Role:         config.RoleGuest,
				LastModified: time.Now(),
				StoragePlan:  param.StoragePlan,
			}
		}
		err = dao.InsertUser(userCol, user)
		if err != nil {
			tools.LogError(err, requestID, "InsertUser failed", user)
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		cred := entity.AccessCredential{
			CloudID:  *config.FlagCloudID,
			UserID:   param.UserID,
			Password: passwd,
		}
		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
			"Data":      []entity.AccessCredential{cred},
		})
		tools.LogInfo("Response PostStoragePlan", requestID, []entity.User{user})
	}
}

func PostMetadata(c *gin.Context) {
	requestID := uuid.New().String()

	var param entity.PostMetadataParam
	err := c.BindJSON(&param)
	if err != nil {
		tools.LogError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	tools.LogInfo("Receive PostMetadata", requestID, param)

	if param.CloudID == *config.FlagCloudID {
		// 来自本云httpserver的请求
		var errs []error
		var routine int
		ch := make(chan error)

		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			tools.LogError(err, requestID, "GetAllCloudInfo failed")
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
			go func(cloud entity.Cloud) {
				err := utils.SendPostMetadata(cloudCol, param, cloud.CloudID)
				if err != nil {
					tools.LogError(err, requestID, "sendPostMetadata failed", param, cloud)
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
			tools.LogError(nil, requestID, "Some sendPostMetadata failed", routine, len(errs))
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
		tools.LogInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	} else {
		// 来自其他云scheduler的请求

		// 校验请求合法性
		_, err = dao.GetCloud(cloudCol, param.CloudID)
		if err != nil {
			tools.LogError(err, requestID, "GetCloudInfo failed", param.CloudID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"RequestID": requestID,
				"Code":      codeUnauthorized,
				"Msg":       errorMsg[codeUnauthorized],
			})
			return
		}

		if param.Type == "Upload" {
			// 写入文件元信息
			err = dao.InsertFiles(fileCol, param.Files)
			if err != nil {
				tools.LogError(err, requestID, "InsertFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// 修改用户存储量
			err = dao.ChangeVolume(userCol, param.UserID, "Upload", param.Files)
			if err != nil {
				tools.LogError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
		} else if param.Type == "Delete" {
			// 删除文件元信息
			err = dao.DeleteFiles(fileCol, param.Files)
			if err != nil {
				tools.LogError(err, requestID, "DeleteFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// 修改用户存储量
			err = dao.ChangeVolume(userCol, param.UserID, "Delete", param.Files)
			if err != nil {
				tools.LogError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
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
			tools.LogError(err, requestID, errorMsg[codeBadRequest])
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
		tools.LogInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	}
}

func Heartbeat(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		<-t.C
		requestID := uuid.New().String()
		tools.LogTrace("Starting to send heartbeat packages", requestID)

		clouds, err := dao.GetOtherClouds(cloudCol, *config.FlagCloudID)
		if err != nil {
			tools.LogError(err, requestID, "GetOtherClouds failed", *config.FlagCloudID)
			continue
		}

		ch := make(chan error)
		param := entity.GetStatusParam{CloudID: *config.FlagCloudID}
		for _, cloud := range clouds {
			go func(cloud entity.Cloud) {
				c, err := utils.SendGetStatus(cloudCol, param, cloud.CloudID)
				if err != nil {
					tools.LogError(err, requestID, "sendGetStatus failed", param, cloud)
					ch <- err
					c = &cloud
					c.Status = "DOWN"
				} else {
					c.Status = "UP"
				}
				err = dao.UpdateCloud(cloudCol, *c)
				if err != nil {
					tools.LogError(err, requestID, "UpdateCloud failed", *c)
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

		tools.LogTrace("Heartbeat finished", requestID, len(clouds), success)
	}
}

func GetAllCloudsStatus(c *gin.Context) {
	requestID := uuid.New().String()

	//查询所有的clouds
	clouds, err := dao.GetAllClouds(cloudCol)
	tools.LogInfo("clouds:", requestID, clouds)
	if err != nil {
		//查询出错，报告错误
		tools.LogError(err, requestID, "query for all clouds status failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}
	//查询成功，返回数据
	//隐藏accesskey和secretkey
	for index := range clouds {
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
	tools.LogInfo("UpdateClouds:", requestID, c)
	var cloud entity.Cloud
	if err := c.ShouldBindJSON(&cloud); err != nil {
		//can't get the clouds
		//return the error
		tools.LogError(err, requestID, "can't get the clouds from paramators")
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}

	//update the clouds
	if err := dao.UpdateCloud(cloudCol, cloud); err != nil {
		//log the error
		tools.LogError(err, requestID, "can't update the cloud")
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	//向其他云同步
	if c.GetHeader("Caller") == "http-server" {
		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			tools.LogError(err, requestID, "can't get other clouds")
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		b, err := json.Marshal(cloud)
		if err != nil {
			tools.LogError(err, requestID, "can't Marshal the cloud")
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, otherCLoud := range clouds {
			if otherCLoud.CloudID != *config.FlagCloudID {
				body := bytes.NewBuffer(b)
				addr := utils.GenAddress(cloudCol, otherCLoud.CloudID, "/update_clouds")
				resp, err := http.Post(addr, "application/json", body)
				if err != nil || resp.StatusCode != 200 {
					tools.LogError(err, requestID, "can't syn to cloud: ", otherCLoud.CloudID)
					//c.JSON(http.StatusBadRequest, gin.H{
					//	"RequestID": requestID,
					//	"Code":      codeInternalError,
					//	"Msg":       errorMsg[codeInternalError],
					//})
					//return
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
	tools.LogInfo("update the clouds succeeded!", requestID, cloud)
}
