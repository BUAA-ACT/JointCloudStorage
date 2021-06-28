package newcloud

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shaoliyin.me/jcspan/dao"
)

const (
	CollectionCloud     = "Cloud"
	CollectionTempCloud	= "tempCloud"
	CollectionVoteCloud = "voteCloud"
	CollectionUser      = "User"
	CollectionFile      = "File"
	MigrationAdvice     = "MigrationAdvice"
	codeOK              = 200
	codeBadRequest      = 400
	codeUnauthorized    = 401
	codeInternalError   = 500
)

var (
	localid string
	errorMsg  = map[int]string{
		codeOK:            "OK",
		codeBadRequest:    "Bad Request",
		codeUnauthorized:  "Unauthorized",
		codeInternalError: "Internal Server Error",
	}
	localMongo                *dao.Dao
	localMongoTempCloud       *dao.Dao
	localMongoVoteRequest *dao.Dao
)


/*
 * NewCloud 的初始化函数，用于初始化mongodb的链接和本地cid
 * mongo：本地mongo数据库地址
 * clouds：database名称
 * cid：本地云的cid
 */
func NewCloudInit(mongo,databasename,cid string) error{
	var err error
	localMongo, err = dao.NewDao(mongo, databasename, CollectionCloud, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		return err
	}

	localMongoTempCloud, err = dao.NewDao(mongo, databasename, CollectionTempCloud, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		return err
	}

	localMongoVoteRequest, err = dao.NewDao(mongo, databasename, CollectionVoteCloud, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		return err
	}

	localid=cid
	return nil
}

func PostNewCloud(c *gin.Context) {
	requestID := uuid.New().String()

	//读取新的cloud
	var tempCloud dao.Cloud
	err := c.ShouldBindJSON(&tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		return
	}

	//2.将新的Cloud存入Mongo中
	temp := dao.VoteCloud{
		Id:      tempCloud.CloudID,
		Cloud:   tempCloud,
		VoteNum: 0,
		Address: tempCloud.Address,
	}
	err = localMongoTempCloud.InsertVoteCloud(temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}
	//3.将新cloud发送给其他的节点,io.reader打包成json返回
	//获取所有云和自己的address
	clouds,err:=localMongo.GetAllClouds()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostNewCloud, message:",err,"RequestID:",requestID)
		return
	}
	for _,cloud:=range clouds{
		if cloud.CloudID==localid{
			temp.Address=cloud.Address
		}
	}
	//将voetcloud发送给其他所有的云
	for _,cloud:=range clouds{
		if cloud.CloudID!=localid{
			b,err:=json.Marshal(clouds)
			if err!=nil{
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				log.Fatal("package:NewCloud, func:PostNewCloud, message:",err," Can't marshal temp vote cloud.","RequestID:",requestID)
				return
			}
			body:=bytes.NewBuffer(b)
			resp,err:=http.Post("http://"+cloud.Address+"/new_cloud_vote","application/json",body)
			if err!=nil{
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				log.Fatal("package:NewCloud, func:PostNewCloud, message:",err," Send temp cloud error! ","RequestID:",requestID)
				return
			}
			log.Info("package:NewCloud, func:PostNewCloud, code:",resp.StatusCode," Code shoule be 200."," RequestID:",requestID)
		}
	}
	c.JSON(http.StatusOK,gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})

}

//接收主节点的投票请求推送，存入VoteRequest表，等待投票
func PostNewCloudVote(c *gin.Context) {
	requestID := uuid.New().String()

	var tempCloud dao.Cloud
	err := c.ShouldBindJSON(&tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostNewCloudRequest, message:",err,"RequestID:",requestID)
		return
	}

	//2.将cloud存入到本地mongo中
	err = localMongoVoteRequest.InsertCloud(tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostNewCloudRequest, message:",err,"RequestID:",requestID)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
	log.Info("package:NewCloud, func:PostNewCloudRequest, RequestID:",requestID,"message:successed!")
}

//获取voteRequest中所有的等待投票的cloud信息
func GetVoteRequest(c * gin.Context){
	requestID:=uuid.New().String()

	//1.读取mongoDB
	clouds,err:=localMongoVoteRequest.GetAllClouds()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:GetVoteRequest, message:",err)
		return
	}

	//2.返回所有clouds
	c.JSON(http.StatusOK,clouds)
	log.Info("package:NewCloud, func:GetVoteRequest, message:successed!")
}

//对一个cloud投票
func PostCloudVote(c * gin.Context){
	requestID:=uuid.New().String()
	var id string
	var req map[string]interface{}

	//获取cloud的id
	err:=c.ShouldBindJSON(&req)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostCloudVote, message:",err,"RequestID:",requestID)
		return
	}
	if _,ok:=req["id"];ok{
		id=req["id"].(string)
	}else{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostCloudVote, message:can't get the id from context","RequestID:",requestID)
		return
	}

	//查询自己的tempcloud，看是不是主节点
	count,err:=localMongoTempCloud.CloudsCount(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostCloudVote, message:",err,"RequestID:",requestID)
		return
	}
	if count>0{
		//若是主节点，则直接加一
		_,err=localMongoTempCloud.AddVoteNum(1,id)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostCloudVote, message:",err,"RequestID:",requestID)
			return
		}
	}else{
		//若不是主节点，给其他云投票
		//获取相关的votecloud信息
		voteCloud,err:=localMongoVoteRequest.GetVoteCloud(id)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostCloudVote, message:",err," Can't get the vote request!"," RequestID:",requestID)
			return
		}

		//序列化投票信息{id：”“，json：num}
		voteMsg:=make(map[string]interface{})
		voteMsg["id"]=voteCloud.Cloud.CloudID
		voteMsg["vote"]=1
		b,err:=json.Marshal(voteMsg)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostCloudVote, message:",err," Can't marshal the cloud","RequestID:",requestID)
			return
		}
		body:=bytes.NewBuffer(b)
		resp,err:=http.Post("http://"+voteCloud.Address+"/master_cloud_vote","application/json",body)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostCloudVote, message:",err," Vote to the master error!"," RequestID:",requestID)
			return
		}
		if resp.StatusCode!=200{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeBadRequest,
				"Msg":       errorMsg[codeBadRequest],
			})
			log.Fatal("package:NewCloud, func:PostCloudVote, message:",err," Vote to the master error!"," RequestID:",requestID)
			return
		}

	}

	c.JSON(http.StatusOK,gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
	log.Info("package:NewCloud, func:GetVoteRequest, message:successed!")
}

//master接收其他云的投票信息
func PostMasterCloudVote(c *gin.Context){
	requestID:=uuid.New().String()
	var id string
	var vote int
	//获取投票id和票数
	var req map[string]interface{}
	err:=c.ShouldBindJSON(&req)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err,"RequestID:",requestID)
		return
	}
	id,ok1:=req["id"].(string)
	vote,ok2:=req["vote"].(int)
	if !(ok1&&ok2){
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:can't get id and vote from context","RequestID:",requestID)
		return
	}

	//向tempcloud表中投票
	modifyNum,err:=localMongoTempCloud.AddVoteNum(vote,id)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err,"RequestID:",requestID)
		return
	}
	if modifyNum<=0{
		//若未查询到相关云，则返回510表示错误
		c.JSON(510,nil)
		return
	}

	//检查现有的投票数量和现有云数量
	voteCloud,err:=localMongoTempCloud.GetVoteCloud(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Can't get vote cloud. ","RequestID:",requestID)
		return
	}

	err=localMongoTempCloud.DeleteVoteCloud(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Delete cloud error. ","RequestID:",requestID)
		return
	}

	totalNum,err:=localMongoTempCloud.GetCloudNum()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Can't get total number. ","RequestID:",requestID)
		return
	}

	//检查是否已经达到多数同步,是则同步到所有的云
	if voteCloud.VoteNum>totalNum/2{
		//将云信息推给其他云
		clouds,err:=localMongo.GetAllClouds()
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Can't get total number. ","RequestID:",requestID)
			return
		}
		var newclouds []dao.Cloud
		newclouds=append(newclouds, voteCloud.Cloud)
		b,err:=json.Marshal(newclouds)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Can't marshal the cloud","RequestID:",requestID)
			return
		}
		body:=bytes.NewBuffer(b)
		for _,cloud:=range clouds{
			_,err:=http.Post("http://"+cloud.Address+"/cloud_syn","application/json",body)
			if err!=nil{
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err,"RequestID:",requestID)
				return
			}
		}

		//向新云同步
		b,err=json.Marshal(clouds)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Can't marshal all clouds","RequestID:",requestID)
			return
		}
		body=bytes.NewBuffer(b)
		_,err=http.Post("http://"+voteCloud.Cloud.Address+"/cloud_syn","application/json",body)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"RequestID": requestID,
				"Code":      codeInternalError,
				"Msg":       errorMsg[codeInternalError],
			})
			log.Fatal("package:NewCloud, func:PostMasterCloudVote, message:",err," Send to new cloud error! ","RequestID:",requestID)
			return
		}
	}
	c.JSON(http.StatusOK,gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
}

//向其他云同步信息
func PostCloudSyn(c *gin.Context){
	requestID:=uuid.New().String()
	var clouds []dao.Cloud

	//获取clouds
	err:=c.ShouldBindJSON(&clouds)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostCloudSyn, message:",err,"RequestID:",requestID)
		return
	}

	//获取现有clouds
	localClouds,err:=localMongo.GetAllClouds()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeBadRequest,
			"Msg":       errorMsg[codeBadRequest],
		})
		log.Fatal("package:NewCloud, func:PostCloudSyn, message:",err,"RequestID:",requestID)
		return
	}

	//将所有的未出现cloud存入mongo
	for _,cloud:=range clouds{
		flag:=true
		for _,localCloud:=range localClouds{
			if localCloud.CloudID==cloud.CloudID{
				flag=false
			}
		}

		if flag{
			err=localMongo.InsertCloud(cloud)
			if err!=nil{
				c.JSON(http.StatusBadRequest, gin.H{
					"RequestID": requestID,
					"Code":      codeInternalError,
					"Msg":       errorMsg[codeInternalError],
				})
				log.Fatal("package:NewCloud, func:PostCloudSyn, message:",err,"RequestID:",requestID)
				return
			}
		}
	}

	c.JSON(http.StatusOK,gin.H{
		"RequestID": requestID,
		"Code":      codeOK,
		"Msg":       errorMsg[codeOK],
	})
}
