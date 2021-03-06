package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"shaoliyin.me/jcspan/utils"
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

type BaseResponse struct {
	RequestID string
	Code      int
	Msg       string
}

type GetStoragePlanParam dao.Preference

type GetStoragePlanData struct {
	StoragePriceFirst dao.StoragePlan
	TrafficPriceFirst dao.StoragePlan
}

type GetDownloadPlanParam struct {
	UserID string
	FileID string
}

type GetDownloadPlanData struct {
	StorageMode string
	Clouds      []dao.Cloud
	Index       []int
}

type GetStatusParam struct {
	CloudID string
}

type GetStatusData struct {
	dao.Cloud
}

type PostStoragePlanParam struct {
	CloudID     string
	UserID      string
	Password    string
	StoragePlan dao.StoragePlan
}

type PostStoragePlanData struct {
	dao.AccessCredential
}

type PostMetadataParam struct {
	CloudID string
	UserID  string
	Type    string
	Clouds  []dao.Cloud
	Files   []dao.File
}

type PostMetadataData struct {
}

func GetStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param GetStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	logInfo("Receive GetStoragePlan", requestID, param)

	clouds, err := db.GetAllClouds()
	if err != nil {
		logError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// ??????????????????
	plan := storagePlan(param, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      plan,
	})

	logInfo("Response GetStoragePlan", requestID, plan)
}

func GetDownloadPlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param GetDownloadPlanParam
	err := c.BindJSON(&param)
	if err != nil {
		logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	logInfo("Receive GetDownloadPlan", requestID, param)

	user, err := db.GetUser(param.UserID)
	if err != nil {
		logError(err, requestID, "GetUserInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	clouds, err := db.GetAllClouds()
	if err != nil {
		logError(err, requestID, "GetAllCloudInfo failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	// ??????????????????
	plan := downloadPlan(user.StoragePlan, clouds)

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      user.StoragePlan,
	})

	logInfo("Response GetDownloadPlan", requestID, plan)
}

func GetStatus(c *gin.Context) {
	requestID := uuid.New().String()
	var param GetStatusParam
	err := c.BindJSON(&param)
	if err != nil {
		logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	logTrace("Receive GetStatus", requestID, param)

	// ??????????????????????????????
	_, err = db.GetCloud(param.CloudID)
	if err != nil {
		logError(err, requestID, "GetCloudInfo failed", param.CloudID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"RequestID": requestID,
			"Code":      codeUnauthorized,
			"Msg":       errorMsg[codeUnauthorized],
		})
		return
	}

	// ??????????????????
	cloud, err := db.GetCloud(*flagCloudID)
	if err != nil {
		logError(err, requestID, "GetCloudInfo failed", *flagCloudID)
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

	logTrace("Response GetStatus", requestID, cloud)
}

func PostStoragePlan(c *gin.Context) {
	requestID := uuid.New().String()

	var param PostStoragePlanParam
	err := c.BindJSON(&param)
	if err != nil {
		logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	logInfo("Receive PostStoragePlan", requestID, param)

	if param.CloudID == *flagCloudID {
		// ????????????httpserver?????????
		plan := &param.StoragePlan
		plan.StoragePrice = calStoragePrice(*plan)
		plan.Availability = calAvailability(*plan)
		plan.TrafficPrice = calTrafficPrice(*plan, false)
		var users []dao.AccessCredential
		ch := make(chan *dao.AccessCredential)

		clouds, err := db.GetAllClouds()
		if err != nil {
			logError(err, requestID, "GetAllCloudInfo failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, cloud := range clouds {
			// ?????????????????????
			if cloud.CloudID == param.CloudID {
				continue
			}
			go func(cloud dao.Cloud) {
				u, err := sendPostStoragePlan(param, cloud.CloudID)
				if err != nil {
					logError(err, requestID, "sendPostStoragePlan failed", param, cloud)
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
			logError(nil, requestID, "Some sendPostStoragePlan failed", len(clouds)-1, len(users))
			//c.JSON(http.StatusInternalServerError, gin.H{
			//	"RequestID": requestID,
			//	"Code":      codeInternalError,
			//	"Msg":       errorMsg[codeInternalError],
			//})
		}

		// ??????????????????
		var user dao.User
		user, err = db.GetUser(param.UserID)
		if err != nil {
			logError(err, requestID, "?????????????????????????????????????????????")
		} else {
			user.StoragePlan = param.StoragePlan
			user.Preference.Vendor = param.StoragePlan.N // ????????????????????????
			user.Preference.StoragePrice = math.Max(user.Preference.StoragePrice, user.StoragePlan.StoragePrice)
			user.Preference.TrafficPrice = math.Max(user.Preference.TrafficPrice, user.StoragePlan.TrafficPrice)
			user.Preference.Availability = math.Min(user.Preference.Availability, user.StoragePlan.Availability)
			err = db.InsertUser(user)
			if err != nil {
				logError(err, requestID, "??????????????????????????????")
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"RequestID": requestID,
			"Code":      codeOK,
			"Msg":       errorMsg[codeOK],
			"Data":      users,
		})
		logInfo("Response PostStoragePlan", requestID, users)
	} else {
		// ???????????????scheduler?????????
		_, err = db.GetCloud(param.CloudID)
		if err != nil {
			logError(err, requestID, "GetCloudInfo failed", param.CloudID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"RequestID": requestID,
				"Code":      codeUnauthorized,
				"Msg":       errorMsg[codeUnauthorized],
			})
			return
		}

		// ????????????
		passwd := param.Password
		// ??????????????????????????????????????? User
		var user dao.User
		if passwd == "" {
			user, err = db.GetUser(param.UserID)
			if err != nil {
				logError(err, requestID, "Get user failed", user)
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
				Password:     AesEncrypt(passwd, *flagAESKey),
				Role:         dao.RoleGuest,
				LastModified: time.Now(),
				StoragePlan:  param.StoragePlan,
			}
		}
		err = db.InsertUser(user)
		if err != nil {
			logError(err, requestID, "InsertUser failed", user)
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
		logInfo("Response PostStoragePlan", requestID, []dao.User{user})
	}
}

func PostMetadata(c *gin.Context) {
	requestID := uuid.New().String()

	var param PostMetadataParam
	err := c.BindJSON(&param)
	if err != nil {
		logError(err, requestID, errorMsg[codeBadRequest])
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}
	logInfo("Receive PostMetadata", requestID, param)

	if param.CloudID == *flagCloudID {
		// ????????????httpserver?????????
		var errs []error
		var routine int
		ch := make(chan error)

		clouds, err := db.GetAllClouds()
		if err != nil {
			logError(err, requestID, "GetAllCloudInfo failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		for _, cloud := range clouds {
			// ???????????????
			if cloud.CloudID == param.CloudID {
				continue
			}
			routine++
			go func(cloud dao.Cloud) {
				err := sendPostMetadata(param, cloud.CloudID)
				if err != nil {
					logError(err, requestID, "sendPostMetadata failed", param, cloud)
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
			logError(nil, requestID, "Some sendPostMetadata failed", routine, len(errs))
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
		logInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	} else {
		// ???????????????scheduler?????????

		// ?????????????????????
		_, err = db.GetCloud(param.CloudID)
		if err != nil {
			logError(err, requestID, "GetCloudInfo failed", param.CloudID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"RequestID": requestID,
				"Code":      codeUnauthorized,
				"Msg":       errorMsg[codeUnauthorized],
			})
			return
		}

		if param.Type == "Upload" {
			// ?????????????????????
			err = db.InsertFiles(param.Files)
			if err != nil {
				logError(err, requestID, "InsertFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// ?????????????????????
			err = db.ChangeVolume(param.UserID, "Upload", param.Files)
			if err != nil {
				logError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
		} else if param.Type == "Delete" {
			// ?????????????????????
			err = db.DeleteFiles(param.Files)
			if err != nil {
				logError(err, requestID, "DeleteFiles failed", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}

			// ?????????????????????
			err = db.ChangeVolume(param.UserID, "Delete", param.Files)
			if err != nil {
				logError(err, requestID, "ChangeVolume failed", param.UserID, "Upload", param.Files)
				c.JSON(http.StatusInternalServerError, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				return
			}
		} else if param.Type == "Migrate" {
			// ?????????????????????PostStoragePlan???????????????????????????????????????
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
			logError(err, requestID, errorMsg[codeBadRequest])
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
		logInfo("Response PostMetadata", requestID, errorMsg[codeOK])
	}
}

func heartbeat(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		<-t.C
		requestID := uuid.New().String()
		logTrace("Starting to send heartbeat packages", requestID)

		clouds, err := db.GetOtherClouds(*flagCloudID)
		if err != nil {
			logError(err, requestID, "GetOtherClouds failed", *flagCloudID)
			continue
		}

		ch := make(chan error)
		param := GetStatusParam{CloudID: *flagCloudID}
		for _, cloud := range clouds {
			go func(cloud dao.Cloud) {
				c, err := sendGetStatus(param, cloud.CloudID)
				if err != nil {
					logError(err, requestID, "sendGetStatus failed", param, cloud)
					ch <- err
					c = &cloud
					c.Status = "DOWN"
				} else {
					c.Status = "UP"
				}
				err = db.UpdateCloud(*c)
				if err != nil {
					logError(err, requestID, "UpdateCloud failed", *c)
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

		logTrace("Heartbeat finished", requestID, len(clouds), success)
	}
}

func GetAllCloudsStatus(c *gin.Context) {
	requestID := uuid.New().String()

	//???????????????clouds
	clouds, err := db.GetAllClouds()
	logInfo("clouds:", requestID, clouds)
	if err != nil {
		//???????????????????????????
		logError(err, requestID, "query for all clouds status failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}
	//???????????????????????????
	//??????accesskey???secretkey
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
	logInfo("UpdateClouds:", requestID, c)
	var cloud dao.Cloud
	if err := c.ShouldBindJSON(&cloud); err != nil {
		//can't get the clouds
		//return the error
		logError(err, requestID, "can't get the clouds from paramators")
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
		logError(err, requestID, "can't update the cloud")
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	//??????????????????
	if c.GetHeader("Caller") == "http-server" {
		clouds, err := db.GetAllClouds()
		if err != nil {
			logError(err, requestID, "can't get other clouds")
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			return
		}

		b, err := json.Marshal(cloud)
		if err != nil {
			logError(err, requestID, "can't Marshal the cloud")
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
				addr := utils.GenAddress(otherCLoud.CloudID, "/update_clouds")
				resp, err := http.Post(addr, "application/json", body)
				if err != nil || resp.StatusCode != 200 {
					logError(err, requestID, "can't syn to cloud: ", otherCLoud.CloudID)
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
	logInfo("update the clouds succeeded!", requestID, cloud)
}
