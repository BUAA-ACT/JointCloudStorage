package controller

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"net/http"
)

type HHH struct {
	A string
	B string
	C string
}

func UserTestPost(con *gin.Context) {
	var a []model.File
	file := model.File{}
	a = append(a, file)
	a = append(a, file)
	a = append(a, file)
	//getQueryAndReturn(con,"EEE")
	requestID, _ := uuid.New()
	var param map[string]interface{}
	con.BindJSON(&param)
	mapHHH := param["HHH"]
	h := HHH{}
	err1 := mapstructure.Decode(mapHHH, &h)
	fmt.Println("exter err1")
	fmt.Println(err1)
	fmt.Println("exter h:")
	fmt.Println(h)
	con.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      args.CodeOK,
	})
}

func UserTestGet(con *gin.Context) {
	requestID, _ := uuid.New()
	var param model.Preference
	err := con.BindJSON(&param)
	fmt.Println("-----")
	fmt.Println(err)
	fmt.Println("-----")

	fmt.Println(param)
	con.JSON(http.StatusOK, gin.H{
		"RequestID": requestID,
		"Code":      args.CodeOK,
		"Data":      gin.H{},
	})
}
