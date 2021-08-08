package utils

import (
	"github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/dao"
	"strings"
)

var (
	addrMap = make(map[string]string)
)

func GenAddress(cloudID, path string) string {
	//if _, ok := addrMap[cloudID]; !ok { // 内存 map 中没有该云节点
	//
	//}
	db := dao.GetDatabaseInstance()
	// Init address map
	clouds, err := db.GetAllClouds()
	if err != nil {
		panic(err)
	}
	var addr string
	for _, c := range clouds {
		if c.CloudID == cloudID {
			addr = CorrectAddress(c.Address)
			break
		}
	}
	if len(addr) == 0 {
		logrus.Errorf("地址生成失败，云节点信息不存在：%v", cloudID)
		return ""
	}
	return "http://" + addr + path
}

func CorrectAddress(addr string) string {
	if !strings.Contains(addr, ":") {
		addr = addr + ":8082"
	}
	return addr
}
