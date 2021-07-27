package controller

import (
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type HHH struct {
	A string
	B string
	C string
}

func UserTestPost(con *gin.Context) {
	fieldRequired := map[string]bool{
		"header": true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)["header"].(string)
	fmt.Println(accessToken)
	return
	//getQueryAndReturnWithHttp(con,"EEE")
	//requestID, _ := uuid.New()
	//var param map[string]interface{}
	//con.BindJSON(&param)
	//mapHHH := param["HHH"]
	//h := HHH{}
	//err1 := mapstructure.Decode(mapHHH, &h)
	//fmt.Println("exter err1")
	//fmt.Println(err1)
	//fmt.Println("exter h:")
	//fmt.Println(h)
	//con.JSON(http.StatusOK, gin.H{
	//	"RequestID": requestID,
	//	"Code":      args.CodeOK,
	//})
}

func UserTestGet(con *gin.Context) {
	fieldRequired := map[string]bool{
		"header": true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)["header"].(string)
	fmt.Println(accessToken)
}

func HeaderTestPost(con *gin.Context) {
	fieldRequired := map[string]bool{
		"header": true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)["header"].(string)
	fmt.Println(accessToken)
}

func CookieTestGet(con *gin.Context) {
	fieldRequired := map[string]bool{
		"gin_cookie1": true,
	}
	valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
		return
	}
	accessToken := (*valueMap)["gin_cookie1"].(string)
	fmt.Println(accessToken)

	//cookie, err := con.Cookie("gin_cookie1")
	//fmt.Println("cookie:" + cookie)
	//if err != nil {
	//	cookie = "NotSet"
	//	con.SetCookie("gin_cookie2", "test", 3600, "/", "localhost", false, true)
	//}
	//fmt.Printf("Cookie value: %s \n", cookie)
}

func WebSocketTestGet(con *gin.Context) {
	ws, weErr := upGrader.Upgrade(con.Writer, con.Request, nil)
	if tools.PrintError(weErr) {
		return
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			tools.PrintError(err)
		}
	}(ws)
	type MessageStruct struct {
		AccessToken string `json:"AccessToken"`
	}
	var a MessageStruct
	// read data from ws
	jsonErr := ws.ReadJSON(&a)
	if jsonErr != nil {
		log.Println("fucking json problem: " + jsonErr.Error())
		return
	}
	go func() {
		for {
			v := gin.H{"message": "fuck you!"}
			writeErr := ws.WriteJSON(v)
			if writeErr != nil {
				log.Println("fucking write json problem: " + jsonErr.Error())
			}
			fmt.Println("hahaha!")
			time.Sleep(1000000000)
		}
	}()

}
