package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	router.GET("/jcspan/*path", router.GetFile)
	rand.Seed(time.Now().Unix())
	return &router
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "JcsPan Transporter")
}

func (router *Router) AddUploadTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	destinationPath := ps.ByName("path")[1:]
	log.Printf("upload to :%v", destinationPath)
	sidCookie, err := r.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie Fail: %v", err)
	}
	// todo: 文件较小时，不落盘，直接内存上传
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// todo: 自定义临时路径
	randStr := genRandomString(10)
	filePath := "./tmp/" + handler.Filename + randStr
	if !isDir("./tmp/") {
		os.MkdirAll("./tmp/", os.ModePerm)
	}
	// 创建文件，且文件必须不存在
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666) // 此处假设当前目录下已存在test目录
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	sourcePath := filePath
	router.processor.CreateTask(model.USER_UPLOAD_SIMPLE, sidCookie.Value, sourcePath, destinationPath)
}

// 获取网盘文件临时下载链接
func (router *Router) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filePath := ps.ByName("path")[1:]
	log.Printf("get tmp download url: %v", filePath)
	sidCookie, err := r.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie File: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Auth Fail")
		return
	}
	task := model.NewTask(0, model.USER_DOWNLOAD_SIMPLE, time.Now(), sidCookie.Value, filePath, "")
	url, err := router.processor.ProcessGetTmpDownloadUrl(task)
	if err != nil {
		log.Printf("Get tmp download url fail: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "500 ERROR")
	}
	fmt.Fprintln(w, url)
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

func genRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	)
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
