package scheduler

import (
	"bytes"
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/code"
	"cloud-storage-httpserver/service/tools"
	"encoding/json"

	"net/http"
)

func GetAllStoragePlanFromScheduler(preference model.GetStoragePlan) (*model.GetStoragePlanResponse, bool) {
	client := http.Client{}
	preferenceJSON, errMarshal := json.Marshal(preference)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/storage_plan", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer resp.Body.Close()
	// debug
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	// decode json
	var response model.GetStoragePlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func GetDownloadPlanFromScheduler(userId string, fileId string) (*model.GetDownloadPlanResponse, bool) {
	client := http.Client{}
	getDownloadPlan := model.GetDownloadPlan{
		UserID: userId,
		FileID: fileId,
	}
	preferenceJSON, errMarshal := json.Marshal(getDownloadPlan)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/download_plan", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer resp.Body.Close()
	// debug
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	// decode json
	var response model.GetDownloadPlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func SetStoragePlanToScheduler(user *model.User, storagePlan *model.StoragePlan) (*model.PostStoragePlanResponse, bool) {
	client := http.Client{}
	postStoragePlan := model.PostStoragePlan{
		CloudID:     *args.CloudID,
		UserID:      user.UserId,
		Password:    code.AesDecrypt(user.Password, *args.EncryptKey),
		StoragePlan: *storagePlan,
	}
	preferenceJSON, errMarshal := json.Marshal(postStoragePlan)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/storage_plan", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer resp.Body.Close()
	// debug
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	// decode json
	var response model.PostStoragePlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}
