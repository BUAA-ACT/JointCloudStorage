package main

import (
	"errors"
	"fmt"
	"shaoliyin.me/jcspan/utils"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"shaoliyin.me/jcspan/dao"
)

type GetStatusResponse struct {
	BaseResponse
	Data dao.Cloud
}

type PostStoragePlanResponse struct {
	BaseResponse
	Data []dao.AccessCredential
}

type PostMetadataResponse struct {
	BaseResponse
}

func sendGetStatus(param GetStatusParam, cloud string) (*dao.Cloud, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetAllowGetMethodPayload(true)

	resp, err := client.R().
		SetBody(param).
		SetResult(&GetStatusResponse{}).
		Get(utils.GenAddress(cloud, "/status"))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return &resp.Result().(*GetStatusResponse).Data, nil
}

func sendPostStoragePlan(param PostStoragePlanParam, cloud string) (*dao.AccessCredential, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetBody(param).
		SetResult(&PostStoragePlanResponse{}).
		Post(utils.GenAddress(cloud, "/storage_plan"))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return &resp.Result().(*PostStoragePlanResponse).Data[0], nil
}

func sendPostMetadata(param PostMetadataParam, cloud string) error {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetBody(param).
		Post(utils.GenAddress(cloud, "/metadata"))

	if err != nil {
		return err
	}
	if resp.StatusCode()/100 != 2 {
		return errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return nil
}

func genAddress(cloudID, path string) string {
	if _, ok := addrMap[cloudID]; !ok { // 内存 map 中没有该云节点
		// Init address map
		clouds, err := db.GetAllClouds()
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
