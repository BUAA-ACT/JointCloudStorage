package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/dao"
	"strings"
)

var (
	addrMap = make(map[string]string)
)

func GenAddress(cloudCol *mongo.Collection, cloudID, path string) string {
	if _, ok := addrMap[cloudID]; !ok { // 内存 map 中没有该云节点
		//db := dao.GetDatabaseInstance()
		// Init address map

		clouds, err := dao.GetAllClouds(cloudCol)
		if err != nil {
			panic(err)
		}
		for _, c := range clouds {
			addrMap[c.CloudID] = c.Address
		}
	}
	addr := addrMap[cloudID]
	if !strings.Contains(addr, ":") {
		addr = addr + ":8082"
	}
	return "http://" + addr + path
}
