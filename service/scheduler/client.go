package main

import (
	"errors"
	"fmt"
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
		Get(genAddress(cloud, "/status"))

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
		Post(genAddress(cloud, "/storage_plan"))

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
		Post(genAddress(cloud, "/metadata"))

	if err != nil {
		return err
	}
	if resp.StatusCode()/100 != 2 {
		return errors.New(fmt.Sprint(resp.StatusCode()))
	}

	return nil
}

func genAddress(cloudID, path string) string {
	addr := addrMap[cloudID]
	if !strings.Contains(addr, ":") {
		addr = addr + ":8082"
	}
	return "http://" + addr + path
}
