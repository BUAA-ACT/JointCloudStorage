package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type RequestTask struct {
	TaskType        string `json:"TaskType"`
	Uid             string `json:"Uid"`
	Sid             string `json:"Sid"`
	DestinationPath string `json:"DestinationPath"`
	StoragePlan     RequestStoragePlan
}

type RequestStoragePlan struct {
	StorageMode string         `json:"StorageMode"`
	Clouds      []RequestCloud `json:"Clouds"`
}

type RequestCloud struct {
	ID string `json:"ID"`
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
	router.GET("/index/*path", router.FileIndex)
	router.POST("/task", router.CreateTask)
	rand.Seed(time.Now().Unix())
	return &router
}

func (router *Router) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var reqTask RequestTask

	err := json.NewDecoder(r.Body).Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch reqTask.TaskType {
	case "Upload":
		var cloudsID []string
		for _, cloud := range reqTask.StoragePlan.Clouds {
			cloudsID = append(cloudsID, cloud.ID)
		}
		task := model.Task{
			Tid:             primitive.NewObjectID(),
			TaskType:        model.UPLOAD,
			State:           model.BLOCKED,
			StartTime:       time.Time{},
			Sid:             reqTask.Sid,
			SourcePath:      "",
			DestinationPath: reqTask.DestinationPath,
			TaskOptions: &model.TaskOptions{
				SourceStoragePlan: nil,
				DestinationPlan: &model.StoragePlan{
					StorageMode: reqTask.StoragePlan.StorageMode,
					Clouds:      cloudsID,
				},
			},
		}
		tid, err := router.processor.taskStorage.AddTask(&task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		tidStr := tid.Hex()
		fmt.Fprintf(w, "%v", tidStr)

	default:
		http.Error(w, "wrong task type", http.StatusNotImplemented)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "JcsPan Transporter")
}

func (router *Router) FileIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := ps.ByName("path")[1:]
	log.Printf("Index path :%v", path)
	sidCookie, err := r.Cookie("sid")
	if err != nil {
		log.Printf("Get sid from cookie Fail: %v", err)
	}
	task := model.NewTask(model.INDEX, time.Now(), sidCookie.Value, path, "")
	for obj := range router.processor.ProcessPathIndex(task) {
		fmt.Fprintf(w, "%s\n", obj.Key)
	}
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
	tid := r.FormValue("tid")
	fmt.Println(tid)
	taskid, err := primitive.ObjectIDFromHex(tid)
	task, err := router.processor.taskStorage.GetTask(taskid)
	if err != nil {
		log.Printf("Get task fail: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// 鉴权
	if sidCookie.Value != task.Sid {
		log.Printf("wrong sid")
		http.Error(w, "wrong sid", http.StatusBadGateway)
		return
	}

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

	task.SourcePath = sourcePath
	task.DestinationPath = destinationPath
	task.State = model.WAITING

	router.processor.taskStorage.SetTask(task.Tid, task)
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
	task := model.NewTask(model.USER_DOWNLOAD_SIMPLE, time.Now(), sidCookie.Value, filePath, "")
	url, err := router.processor.ProcessGetTmpDownloadUrl(task)
	if err != nil {
		log.Printf("Get tmp download url fail: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "500 ERROR")
	}
	fmt.Fprintln(w, url)
}

func (router *Router) SimpleSync(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	srcPath := r.FormValue("srcpath")
	dstPath := r.FormValue("dstpath")
	sid := r.FormValue("sid")
	router.processor.CreateTask(model.SYNC_SIMPLE, sid, srcPath, dstPath)
	fmt.Println("task simple sync created success")
}

func (router *Router) SimpleUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
