package scheduler

import (
	"bytes"
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"encoding/json"
	"io"

	"net/http"
)

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	tools.PrintError(err)
}

func GetAllStoragePlanFromScheduler(preference *model.GetStoragePlan) (*model.GetStoragePlanResponse, bool) {
	client := http.Client{}
	preferenceJSON, errMarshal := json.Marshal(*preference)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/storage_plan", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.GetStoragePlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func GetDownloadPlanFromScheduler(userID string, fileId string) (*model.GetDownloadPlanResponse, bool) {
	client := http.Client{}
	getDownloadPlan := model.GetDownloadPlan{
		UserID: userID,
		FileID: fileId,
	}
	getDownloadPlanJSON, errMarshal := json.Marshal(getDownloadPlan)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/download_plan", bytes.NewReader(getDownloadPlanJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.GetDownloadPlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func SetStoragePlanToScheduler(userID string, storagePlan *model.StoragePlan) (*model.PostStoragePlanResponse, bool) {
	client := http.Client{}
	postStoragePlan := model.PostStoragePlan{
		CloudID:     *args.CloudID,
		UserID:      userID,
		StoragePlan: *storagePlan,
	}
	postStoragePlanJson, errMarshal := json.Marshal(postStoragePlan)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/storage_plan", bytes.NewReader(postStoragePlanJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.PostStoragePlanResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func GetAllCloudsFromScheduler() (*model.GetAllCloudsResponse, bool) {
	client := http.Client{}
	getAllClouds := model.GetAllClouds{}
	getAllCloudsJSON, errMarshal := json.Marshal(getAllClouds)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/all_clouds_status", bytes.NewReader(getAllCloudsJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.GetAllCloudsResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	if response.Data == nil {
		response.Data = make([]model.Cloud, 0)
	}
	return &response, true
}

func PostUpdateCloudToScheduler(updateCloud *model.PostUpdateCloud) (*model.PostUpdateCloudResponse, bool) {
	client := http.Client{}
	postUpdateCloudJson, errMarshal := json.Marshal(*updateCloud)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/update_clouds", bytes.NewReader(postUpdateCloudJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.PostUpdateCloudResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func PostNewCloudToScheduler(newCloud *model.PostNewCloud) (*model.PostNewCloudResponse, bool) {
	client := http.Client{}
	postNewCloudJson, errMarshal := json.Marshal(*newCloud)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/new_cloud", bytes.NewReader(postNewCloudJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.PostNewCloudResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func GetVoteRequestsFromScheduler() (*model.GetVoteRequestsResponse, bool) {
	client := http.Client{}
	getVoteRequests := model.GetVoteRequests{}
	getVoteRequestsJson, errMarshal := json.Marshal(getVoteRequests)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("GET", *args.SchedulerUrl+"/vote_request", bytes.NewReader(getVoteRequestsJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.GetVoteRequestsResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	if response.Data == nil {
		response.Data = make([]model.Cloud, 0)
	}
	return &response, true
}

func PostCloudVoteToScheduler(cloudID string, voteResult bool) (*model.PostCloudVoteResponse, bool) {
	client := http.Client{}
	postCloudVote := model.PostCloudVote{
		CloudID:    cloudID,
		VoteResult: voteResult,
	}
	postCloudVoteJson, errMarshal := json.Marshal(postCloudVote)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/cloud_vote", bytes.NewReader(postCloudVoteJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.PostCloudVoteResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func PostKeyToScheduler(key *model.AccessKey) (*model.PostKeyToSchedulerResponse, bool) {
	client := http.Client{}
	postKeyJson, errMarshal := json.Marshal(*key)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/add_key", bytes.NewReader(postKeyJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.PostKeyToSchedulerResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func DeleteKeyToScheduler(userID string, accessKey string) (*model.DeleteKeyToSchedulerResponse, bool) {
	client := http.Client{}
	deleteKey := model.AccessKey{
		UserID:    userID,
		AccessKey: accessKey,
	}
	deleteKeyJson, errMarshal := json.Marshal(deleteKey)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	req, errNewRequest := http.NewRequest("POST", *args.SchedulerUrl+"/delete_key", bytes.NewReader(deleteKeyJson))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	req.Header.Set(args.HttpHeaderKeyForScheduler, args.HttpHeaderValueMe)
	resp, errDoRequest := client.Do(req)
	if tools.PrintError(errDoRequest) {
		return nil, false
	}
	defer closeBody(resp.Body)
	// decode json
	var response model.DeleteKeyToSchedulerResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}
