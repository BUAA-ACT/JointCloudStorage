package newcloud

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"shaoliyin.me/jcspan/entity"
	"shaoliyin.me/jcspan/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shaoliyin.me/jcspan/dao"
)

func PostNewCloud(c *gin.Context) {
	requestID := uuid.New().String()

	//读取新的cloud
	var tempCloud entity.Cloud
	err := c.ShouldBindJSON(&tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      "Get new cloud error:" + err.Error(),
		})
		return
	}

	//2.将新的Cloud存入Mongo中的tempCloud和voteCLoud
	temp := entity.VoteCloud{
		Id:      tempCloud.CloudID,
		Cloud:   tempCloud,
		VoteNum: 0,
		Address: tempCloud.Address,
	}
	err = dao.InsertVoteCloud(tempCloudCol, temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "New cloud store tempCloud error:" + err.Error(),
		})
		return
	}

	//3.将新cloud发送给其他的节点
	//获取所有云信息
	clouds, err := dao.GetAllClouds(cloudCol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "get all clouds error:" + err.Error(),
		})
		log.Error("get all clouds err package:NewCloud, func:PostNewCloud, message:", err, "RequestID:", requestID)
		return
	}
	//初始化tempCloud，将address设为本地的address
	for _, cloud := range clouds {
		if cloud.CloudID == localID {
			temp.Address = cloud.Address
		}
	}
	//插入到本地的VoteCLoud
	err = dao.InsertVoteCloud(voteCloudCol, temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "New cloud store voteCloud error:" + err.Error(),
		})
		return
	}
	//并将tempCloud序列化，生成io.reader
	b, err := json.Marshal(temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "tempCloud marshal error:" + err.Error(),
		})
		log.Error("gen json err, package:NewCloud, func:PostNewCloud, message:", err, " Can't marshal temp vote cloud.", "RequestID:", requestID)
		return
	}
	//将voetcloud发送给其他所有的云
	//本地测试不测试本部分内容
	if env != "localDebug" {
		for _, cloud := range clouds {
			if cloud.CloudID != localID {
				body := bytes.NewBuffer(b)
				addr := utils.GenAddress(cloudCol, cloud.CloudID, "/new_cloud_vote")
				resp, err := http.Post(addr, "application/json", body)
				if err != nil {
					//c.JSON(http.StatusBadRequest, gin.H{
					//	"RequestID": requestID,
					//	"Code":      codeInternalError,
					//	"Msg":       errorMsg[codeInternalError],
					//	"Test":      "Send to other clouds error:" + err.Error(),
					//})
					log.Error("send to other clouds error, package:NewCloud, func:PostNewCloud, message:", err, " Send temp cloud error! ", "RequestID:", requestID)
					//return
				}
				log.Info("package:NewCloud, func:PostNewCloud, code:", resp.StatusCode, " Code shoule be 200.", " RequestID:", requestID)
			}
		}
	} else if env == "localDebug" {
		log.Info("This is local debug, don't send voteCloud to other clouds.\n")
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})

}

// PostNewCloudVote 接收主节点的投票请求推送，存入VoteRequest表，等待投票
func PostNewCloudVote(c *gin.Context) {
	requestID := uuid.New().String()

	var tempCloud entity.VoteCloud
	err := c.ShouldBindJSON(&tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      err,
		})
		log.Error("Bind Json error package:NewCloud, func:PostNewCloudRequest, message:", err, "RequestID:", requestID)
		log.Error(c.Request.Body)
		return
	}

	//2.将cloud存入到本地mongo中
	err = dao.InsertVoteCloud(voteCloudCol, tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      err,
		})
		log.Error("insert vote Cloud error package:NewCloud, func:PostNewCloudRequest, message:", err, "RequestID:", requestID)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Test":      "",
	})
	log.Info("package:NewCloud, func:PostNewCloudRequest, RequestID:", requestID, "message:successed!")
}

// GetVoteRequest 获取voteRequest中所有的等待投票的cloud信息
func GetVoteRequest(c *gin.Context) {
	requestID := uuid.New().String()

	//1.读取mongoDB
	clouds, err := dao.GetAllVoteCloud(voteCloudCol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "Read from mongo error:" + err.Error(),
		})
		log.Error("get all vote cloud from db error, package:NewCloud, func:GetVoteRequest, message:", err)
		return
	}

	//解析voteCloud，获取cloud信息
	var cloudMsg []entity.Cloud
	for _, cloud := range clouds {
		cloudMsg = append(cloudMsg, cloud.Cloud)
	}

	//2.返回所有clouds
	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
		"Data":      cloudMsg,
	})
	log.Info("package:NewCloud, func:GetVoteRequest, message:successed!")
}

// PostCloudVote 对一个cloud投票
func PostCloudVote(c *gin.Context) {
	requestID := uuid.New().String()
	var id string
	var req map[string]interface{}

	//获取cloud的id
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Error("bind json err, package:NewCloud, func:PostCloudVote, message:", err, " Can't get id.", "RequestID:", requestID)
		return
	}
	if _, ok := req["CloudID"]; ok {
		id = req["CloudID"].(string)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Error("package:NewCloud, func:PostCloudVote, message:can't get the id from context", "RequestID:", requestID)
		return
	}

	//查询自己的tempcloud，看是不是主节点
	count, err := dao.CloudsCount(tempCloudCol, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Error("count clouds err, package:NewCloud, func:PostCloudVote, message:", err, " count cloud error.", "RequestID:", requestID)
		return
	}
	if count > 0 {
		//若是主节点，则直接加一,并检查是否查过半数
		err = voteCheck(id, 1)
		if err == tempNotFound {
			c.JSON(510, nil)
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
				"Test":      "voteCheck fail:" + err.Error(),
			})
			return
		}
	} else {
		//若不是主节点，给其他云投票
		//获取相关的votecloud信息
		voteCloud, err := dao.GetVoteCloud(voteCloudCol, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Error("get vote cloud err, package:NewCloud, func:PostCloudVote, message:", err, " Can't get the vote request!", " RequestID:", requestID)
			return
		}

		//序列化投票信息{id：”“，json：num}
		voteMsg := make(map[string]interface{})
		voteMsg["CloudID"] = voteCloud.Cloud.CloudID
		voteMsg["vote"] = 1
		b, err := json.Marshal(voteMsg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Error("marshal voteMsg err, package:NewCloud, func:PostCloudVote, message:", err, " Can't marshal the cloud", "RequestID:", requestID)
			return
		}
		//发出投票请求
		body := bytes.NewBuffer(b)
		if env != "localDebug" {
			addr := "http://" + utils.CorrectAddress(voteCloud.Address) + "/master_cloud_vote"
			resp, err := http.Post(addr, "application/json", body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				log.Error("post vote msg err, package:NewCloud, func:PostCloudVote, message:", err, " Vote to the master error!", " RequestID:", requestID)
				return
			}
			//处理返回结果
			if resp.StatusCode == 510 {
				c.JSON(510, nil)
			} else if resp.StatusCode != 200 {
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeBadRequest,
					"Msg":       errorMsg[codeBadRequest],
				})
				log.Error("vote msg statusCode err, package:NewCloud, func:PostCloudVote, message:", err, " Vote to the master error!", " RequestID:", requestID)
				return
			}
		}

	}
	//删除voteCloud
	err = dao.DeleteVoteCloud(voteCloudCol, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
			"Test":      "Delete cloud fail:" + err.Error(),
		})
		log.Error("delete vote cloud err, package:NewCloud, func:PostCloudVote, message:", err, " Delete cloud error. ", "RequestID:", requestID)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
	log.Info("package:NewCloud, func:GetVoteRequest, message:successed!")
}

// PostMasterCloudVote master接收其他云的投票信息
func PostMasterCloudVote(c *gin.Context) {
	requestID := uuid.New().String()
	var id string
	var vote int
	//获取投票id和票数
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test:":     "Get vote request fail:" + err.Error(),
		})
		log.Error("json 解析失败 package:NewCloud, func:PostMasterCloudVote, message:", err, "RequestID:", requestID)
		return
	}
	_, ok1 := req["CloudID"]
	_, ok2 := req["vote"]
	if !(ok1 && ok2) {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      "Get field id and vote fail.",
		})
		log.Error("未获取到 id 和 vote， package:NewCloud, func:PostMasterCloudVote, message:Can't get id and vote from context\n",
			req, "RequestID:", requestID)
		return
	}
	id = req["CloudID"].(string)
	vote = int(req["vote"].(float64))

	//处理投票
	err = voteCheck(id, vote)
	if err == tempNotFound {
		c.JSON(510, nil)
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      "voteCheck fail:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
}

//获取投票新云的id和投票数，对新云进行投票并检查是否超过半数
func voteCheck(id string, vote int) error {
	modifyNum, err := dao.AddVoteNum(tempCloudCol, vote, id)
	if err != nil {
		log.Error("投票记录入库失败， package:NewCloud, func:voteCheck, message:", err)
		return err
	}
	if modifyNum <= 0 {
		//若未查询到相关云，则返回510表示错误
		return tempNotFound
	}

	//检查现有的投票数量和现有云数量
	voteCloud, err := dao.GetVoteCloud(tempCloudCol, id)
	if err != nil {
		log.Error("获取现有投票数量失败，package:NewCloud, func:voteCheck, message:", err, " Can't get vote cloud. ")
		return err
	}

	//获取所有云的数量
	totalNum, err := dao.GetCloudNum(cloudCol)
	if err != nil {
		log.Error("获取所有云信息失败，package:NewCloud, func:voteCheck, message:", err, " Can't get total number. ")
		return err
	}

	//检查是否已经达到多数同步,是则同步到所有的云
	if voteCloud.VoteNum > totalNum/2 {
		//获取所有云信息
		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			log.Error("获取所有云信息失败，package:NewCloud, func:voteCheck, message:", err, " Can't get all clouds. ")
			return err
		}

		//写入本地Cloud表
		err = dao.InsertCloud(cloudCol, voteCloud.Cloud)
		if err != nil {
			log.Error("写入本地 cloud 表失败，package:NewCloud, func:voteCheck, message:", err)
			return err
		}

		//封装新云信息
		var newclouds []entity.Cloud
		newclouds = append(newclouds, voteCloud.Cloud)
		b, err := json.Marshal(newclouds)
		if err != nil {
			log.Error("解析新云信息失败，package:NewCloud, func:voteCheck, message:", err, " Can't marshal the cloud")
			return err
		}
		var body *bytes.Buffer
		//同步云信息
		for _, cloud := range clouds {
			if cloud.CloudID != localID && env != "localDebug" {
				body = bytes.NewBuffer(b)
				addr := utils.GenAddress(cloudCol, cloud.CloudID, "/cloud_syn")
				_, err := http.Post(addr, "application/json", body)
				if err != nil {
					log.Error("发送新云信息到其他节点失败，package:NewCloud, func:voteCheck, message:", err)
					return err
				}
			}

		}

		//向新云同步 所有云节点信息？？
		b, err = json.Marshal(clouds)
		if err != nil {
			log.Error("package:NewCloud, func:PostMasterCloudVote, message:", err, " Can't marshal all clouds")
			return err
		}

		body = bytes.NewBuffer(b)
		addr := utils.GenAddress(cloudCol, voteCloud.Cloud.CloudID, "/cloud_syn")
		_, err = http.Post(addr, "application/json", body)
		if err != nil {
			log.Error("向新云同步所有云节点信息失败，package:NewCloud, func:PostMasterCloudVote, message:", err, " Send to new cloud error! ", "voteNum:", voteCloud, "totalNum:", totalNum)
			return err
		}

		//删除tempCloud
		err = dao.DeleteVoteCloud(tempCloudCol, id)
		if err != nil {
			log.Error("删除临时云数据失败，package:NewCloud, func:PostMasterCloudVote, message:", err, " Delete cloud error. ")
			return err
		}
	}
	return nil
}

// PostCloudSyn 向其他云同步信息
//接收一个cloud的数组，并与已有的clouds进行对比
//若id不同，则将其存入collection Cloud里
func PostCloudSyn(c *gin.Context) {
	requestID := uuid.New().String()
	var clouds []entity.Cloud

	//获取clouds
	err := c.ShouldBindJSON(&clouds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      "Get new cloud error:" + err.Error(),
		})
		log.Error("package:NewCloud, func:PostCloudSyn, message:", err, "RequestID:", requestID)
		return
	}

	//获取现有clouds
	localClouds, err := dao.GetAllClouds(cloudCol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
			"Test":      "Read Collection Clouds error:" + err.Error(),
		})
		log.Error("获取现有云信息失败，package:NewCloud, func:PostCloudSyn, message:", err, "RequestID:", requestID)
		return
	}

	//将所有的未出现cloud存入mongo
	for _, cloud := range clouds {
		flag := true
		for _, localCloud := range localClouds {
			if localCloud.CloudID == cloud.CloudID {
				flag = false
			}
		}

		if flag {
			err = dao.InsertCloud(cloudCol, cloud)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
					"Test":      "Insert new cloud to collection Cloud error:" + err.Error(),
				})
				log.Error("新云信息插入数据库失败 package:NewCloud, func:PostCloudSyn, message:", err, "RequestID:", requestID)
				return
			}
		}

		err = dao.DeleteVoteCloud(voteCloudCol, cloud.CloudID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
				"Test":      "Insert new cloud to collection Cloud error:" + err.Error(),
			})
			log.Error("新云信息删除voteCloud失败 package:NewCloud, func:PostCloudSyn, message:", err, "RequestID:", requestID)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
}
