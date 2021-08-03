package utils

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/dao"
	"shaoliyin.me/jcspan/entity"
	"strings"
	"time"
)

var (
	addrMap = make(map[string]string)
)

func SendGetStatus(cloudCol *mongo.Collection, param entity.GetStatusParam, cloud string) (*entity.Cloud, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetAllowGetMethodPayload(true)

	resp, err := client.R().
		SetBody(param).
		SetResult(&entity.GetStatusResponse{}).
		Get(GenAddress(cloudCol, cloud, "/status"))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return &resp.Result().(*entity.GetStatusResponse).Data, nil
}

func SendPostStoragePlan(cloudCol *mongo.Collection, param entity.PostStoragePlanParam, cloud string) (*entity.AccessCredential, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetBody(param).
		SetResult(&entity.PostStoragePlanResponse{}).
		Post(GenAddress(cloudCol, cloud, "/storage_plan"))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return &resp.Result().(*entity.PostStoragePlanResponse).Data[0], nil
}

func SendPostMetadata(cloudCol *mongo.Collection, param entity.PostMetadataParam, cloud string) error {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetBody(param).
		Post(GenAddress(cloudCol, cloud, "/metadata"))

	if err != nil {
		return err
	}
	if resp.StatusCode()/100 != 2 {
		return errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return nil
}

//func genAddress(cloudID, path string) string {
//	if _, ok := addrMap[cloudID]; !ok { // 内存 map 中没有该云节点
//		// Init address map
//		clouds, err := db.GetAllClouds()
//		if err != nil {
//			panic(err)
//		}
//		for _, c := range clouds {
//			addrMap[c.CloudID] = c.Address
//		}
//	}
//	addr := addrMap[cloudID]
//	if !strings.Contains(addr, ":") {
//		addr = addr + ":8082"
//	}
//	return "http://" + addr + path
//}

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
