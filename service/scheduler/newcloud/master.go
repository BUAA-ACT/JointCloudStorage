package newcloud

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shaoliyin.me/jcspan/dao"
)

const (
	Version = "v0.2"

	CollectionCloud     = "Cloud"
	CollectionUser      = "User"
	CollectionFile      = "File"
	MigrationAdvice     = "MigrationAdvice"
	TempCloud           = "TempCloud"
	TempCloudFromOthers = "TempCloudFromOthers"
	codeOK              = 200
	codeBadRequest      = 400
	codeUnauthorized    = 401
	codeInternalError   = 500
)

var (
	flagMongo = flag.String("mongo", "mongodb://localhost:27017", "mongodb address")
	flagEnv   = flag.String("env", "test", "dev|test|prod")
	errorMsg  = map[int]string{
		codeOK:            "OK",
		codeBadRequest:    "Bad Request",
		codeUnauthorized:  "Unauthorized",
		codeInternalError: "Internal Server Error",
	}
	localMongo                *dao.Dao
	localMongoTempCloud       *dao.Dao
	localMongoCloudFromOthers *dao.Dao
)

type localCloud struct {
	id      string
	cloud   dao.Cloud
	voteNum int
}

/*
 *init func，create the connection to local mongo
 */
func init() {
	var err error
	localMongo, err = dao.NewDao(*flagMongo, *flagEnv, CollectionCloud, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		panic(err)
	}

	localMongoTempCloud, err = dao.NewDao(*flagMongo, *flagEnv, TempCloud, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		panic(err)
	}

	localMongoCloudFromOthers, err = dao.NewDao(*flagMongo, *flagEnv, TempCloudFromOthers, CollectionUser, CollectionFile, MigrationAdvice)
	if err != nil {
		panic(err)
	}
}

func PostNewCloud(c *gin.Context) {
	requestID := uuid.New().String()

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
	temp := localCloud{
		id:      tempCloud.CloudID,
		cloud:   tempCloud,
		voteNum: 1,
	}
	err = localMongoTempCloud.InsertTempCloud(temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}
	//3.将新cloud发送给其他的节点,io.reader打包成json返回
	//data,err:=json.Marshal(tempCloud)
	//request,_:=http.NewRequest("POST","",bytes.NewBuffer(data))

}

func PostCloudFromOthers(c *gin.Context) {
	requestID := uuid.New().String()

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

	//2.将cloud存入到本地mongo中
	err = localMongoCloudFromOthers.InsertCloud(tempCloud)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"RequestID": requestID,
			"Code":      codeInternalError,
			"Msg":       errorMsg[codeInternalError],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      codeInternalError,
		"Msg":       errorMsg[codeInternalError],
	})
}
