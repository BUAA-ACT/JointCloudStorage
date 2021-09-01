package keySyn

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/entity"
	"shaoliyin.me/jcspan/utils"
)

func PostKeyUpsert(c *gin.Context) {
	requestId := uuid.New().String()

	//获取ak
	var ak entity.AccessKey
	err := c.ShouldBindJSON(&ak)
	if err != nil {
		log.Error("package:keySyn, func:PostKeyUpsert,获取accesskey失败:", err, "requestID:", requestId)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Code": 500,
			"Test": "can't get accesskey from context" + err.Error(),
		})
	}

	//key同步
	err = keySyn(ak, c.GetHeader(CallerHeaderName), SynTypeUpsert)
	if err != nil {
		log.Error("package:keySyn, func:keySyn", err, "requestID:", requestId)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Code":      500,
			"ReuqestID": requestId,
			"Test":      "key syn failed, " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":      http.StatusOK,
		"RequestID": requestId,
	})
}

func keySyn(ak entity.AccessKey, caller string, synType string) error {
	var err error
	switch synType {
	case SynTypeUpsert:
		err = dao.KeyUpsert(keyCol, ak)
	case SynTypeDelete:
		err = dao.DeleteKey(keyCol, ak)
	}

	if err != nil {
		return errors.New("同步本地key失败，" + err.Error())
	}

	if caller == CallerHttpServer {
		//若调用者是http-server,向其他节点同步
		//获取其他clouds
		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			return errors.New("获取clouds失败，" + err.Error())
		}

		//序列化
		b, err := json.Marshal(ak)
		if err != nil {
			return errors.New("序列化失败，" + err.Error())
		}
		//开始同步操作
		for _, cloud := range clouds {
			body := bytes.NewBuffer(b)
			if cloud.CloudID != localCid {
				var req *http.Request
				switch synType {
				case SynTypeUpsert:
					addr := utils.GenAddress(cloudCol, cloud.CloudID, endPointAddKey)
					req, err = http.NewRequest("POST", addr, body)
				case SynTypeDelete:
					addr := utils.GenAddress(cloudCol, cloud.CloudID, endPointDeleteKey)
					req, err = http.NewRequest("POST", addr, body)
				}

				if err != nil {
					return errors.New("生成NewRequest失败," + err.Error())
				}
				req.Header.Set(CallerHeaderName, CallerScheduler)

				//发送请求
				client := http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Error("发送key同步http请求失败，错误详情：", err.Error())
					continue
				} else if resp.StatusCode != 200 {
					log.Error("key同步操作失败，Code：", resp.StatusCode, ",错误详情：", resp.StatusCode)
				}
			}
		}
	}
	return nil
}

func PostKeyDelete(c *gin.Context) {
	requestId := uuid.New().String()

	//获取ak
	var ak entity.AccessKey
	err := c.ShouldBindJSON(&ak)
	if err != nil {
		log.Error("package:keySyn, func:PostDeleteKey,获取accesskey失败:", err, "requestID:", requestId)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Code": 500,
			"Test": "can't get accesskey from context" + err.Error(),
		})
	}

	//key删除
	err = keySyn(ak, c.GetHeader(CallerHeaderName), SynTypeDelete)
	if err != nil {
		log.Error("package:keySyn, func:keySyn", err, "requestID:", requestId)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Code":      500,
			"ReuqestID": requestId,
			"Test":      "key syn failed, " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":      http.StatusOK,
		"RequestID": requestId,
	})
}
