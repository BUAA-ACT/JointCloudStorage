package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Cors(con *gin.Context) {
	method := con.Request.Method
	//if origin != "" {
	con.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
	con.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	con.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	con.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	con.Header("Access-Control-Allow-Credentials", "true")
	//}
	if method == "OPTIONS" {
		con.AbortWithStatus(http.StatusNoContent)
	}
	con.Next()
}
