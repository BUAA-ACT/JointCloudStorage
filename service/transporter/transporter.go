package transporter

import (
	"log"
	"net/http"
)

// 开启 transporter 服务
func StartServe() {
	router := NewRouter()
	log.Println(http.ListenAndServe(":9648", router))
}
