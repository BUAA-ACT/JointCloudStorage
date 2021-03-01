package Transporter

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"fmt"
	"log"
	"net/http"
)

func main() {
	StartServe()
}

// 开启 transporter 服务
func StartServe() {
	// 初始化任务数据库
	storage := model.NewInMemoryTaskStorage()
	processor := controller.TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(model.NewSimpleInMemoryStorageDatabase())
	// 初始化路由
	router := controller.NewRouter(processor)
	fmt.Println("Transporter Started")
	log.Println(http.ListenAndServe(":9648", router))
}
