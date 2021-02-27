package transporter

import (
	"log"
	"net/http"
)

// 开启 transporter 服务
func StartServe() {
	// 初始化任务数据库
	storage := NewInMemoryTaskStorage()
	processor := TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(NewSimpleInMemoryStorageDatabase())
	// 初始化路由
	router := NewRouter(processor)

	log.Println(http.ListenAndServe(":9648", router))
}
