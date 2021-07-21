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

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

// Start is to start a ws server
//func (manager *ClientManager) Start() {
//	for {
//		select {
//		case conn := <-manager.Register:
//			manager.Clients[conn] = true
//			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
//			manager.Send(jsonMessage, conn)
//		case conn := <-manager.Unregister:
//			if _, ok := manager.Clients[conn]; ok {
//				close(conn.Send)
//				delete(manager.Clients, conn)
//				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
//				manager.Send(jsonMessage, conn)
//			}
//		case message := <-manager.Broadcast:
//			for conn := range manager.Clients {
//				select {
//				case conn.Send <- message:
//				default:
//					close(conn.Send)
//					delete(manager.Clients, conn)
//				}
//			}
//		}
//	}
//}

// Send is to send ws message to ws client
//func (manager *ClientManager) Send(message []byte, ignore *Client) {
//	for conn := range manager.Clients {
//		if conn != ignore {
//			conn.Send <- message
//		}
//	}
//}

//func (c *Client) Read() {
//	defer func() {
//		Manager.Unregister <- c
//		c.Socket.Close()
//	}()
//
//	for {
//		_, message, err := c.Socket.ReadMessage()
//		if err != nil {
//			Manager.Unregister <- c
//			c.Socket.Close()
//			break
//		}
//		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
//		Manager.Broadcast <- jsonMessage
//	}
//}
//
//func (c *Client) Write() {
//	defer func() {
//		c.Socket.Close()
//	}()
//
//	for {
//		select {
//		case message, ok := <-c.Send:
//			if !ok {
//				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			c.Socket.WriteMessage(websocket.TextMessage, message)
//		}
//	}
//}

func WebSocketTestGet(con *gin.Context) {
	ws, weErr := upGrader.Upgrade(con.Writer, con.Request, nil)
	if tools.PrintError(weErr) {
		return
	}
	//if err != nil {
	//	return
	//}
	//defer func(ws *websocket.Conn) {
	//	err := ws.Close()
	//	if err != nil {
	//		tools.PrintError(err)
	//	}
	//}(ws)
	//fieldRequired := map[string]bool{
	//	args.FieldWordAccessToken: true,
	//}
	//valueMap, existMap := getQueryAndReturnWithHttp(con, &fieldRequired)
	//if tools.RequiredFieldNotExist(&fieldRequired, existMap) {
	//	return
	//}
	//accessToken := (*valueMap)[args.FieldWordAccessToken].(string)
	////check token
	//_, _, valid := UserCheckAccessToken(con, accessToken, &[]string{args.UserAllRole})
	//if !valid {
	//	return
	//}

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
