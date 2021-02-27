package transporter

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type Router struct {
	*httprouter.Router
	processor TaskProcessor
}

func NewRouter(processor TaskProcessor) *Router {
	var router Router
	router = Router{
		Router:    httprouter.New(),
		processor: processor,
	}
	router.GET("/", Index)
	router.POST("/upload/*path", router.AddUploadTask)
	return &router
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "JcsPan Transporter")
}

func (router *Router) AddUploadTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	destinationPath := ps.ByName("path")
	log.Printf("upload to :%v", destinationPath)
	sidCookie, err := r.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie Fail: %v", err)
	}
	sourcePath := ""
	router.processor.CreateTask(USER_UPLOAD_SIMPLE, sidCookie.Value, sourcePath, destinationPath)
}

func NewTestRouter(processor TaskProcessor) *Router {
	router := NewRouter(processor)
	router.GET("/test/setcookie", testSetCookie)
	return router
}

func testSetCookie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	expiration := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:    "sid",
		Value:   "tttteeeesssstttt",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, "Cookie Already Set")
}
