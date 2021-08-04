package utils

import (
	"shaoliyin.me/jcspan/dao"
	"strings"
)

var (
	addrMap = make(map[string]string)
)

func GenAddress(cloudID, path string) string {
	if _, ok := addrMap[cloudID]; !ok { // 内存 map 中没有该云节点
		db := dao.GetDatabaseInstance()
		// Init address map
		clouds, err := db.GetAllClouds()
		if err != nil {
			panic(err)
		}
		for _, c := range clouds {
			addrMap[c.CloudID] = c.Address
		}
	}
	addr := CorrectAddress(addrMap[cloudID])
	return "http://" + addr + path
}

func CorrectAddress(addr string)string{
	if !strings.Contains(addr,":"){
		addr=addr+":8082"
	}
	return addr
}
