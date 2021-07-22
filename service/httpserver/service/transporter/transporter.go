package transporter

import (
	"bytes"
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"encoding/json"
	"net/http"
)

func PreUploadFile(path string, user *model.User) (*model.TaskResponse, bool) {
	task := model.TaskRequest{
		TaskType:               args.TaskTypeUpload,
		UserID:                 user.UserID,
		DestinationPath:        path,
		DestinationStoragePlan: user.StoragePlan,
	}
	preferenceJSON, errMarshal := json.Marshal(task)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	resp, errNewRequest := http.Post(*args.TransporterUrl+"/task", "application/json", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	var response model.TaskResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func DownLoadFile(downloadName string, userID string, downloadPlan model.StoragePlan) (*model.TaskResponse, bool) {
	task := model.TaskRequest{
		TaskType:          args.TaskTypeDownload,
		UserID:            userID,
		SourcePath:        downloadName,
		SourceStoragePlan: downloadPlan,
	}
	preferenceJSON, errMarshal := json.Marshal(task)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	resp, errNewRequest := http.Post(*args.TransporterUrl+"/task", "application/json", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	var response model.TaskResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}

func DeleteFile(deleteName string, user *model.User) bool {
	task := model.TaskRequest{
		TaskType:          args.TaskTypeDelete,
		UserID:            user.UserID,
		SourcePath:        deleteName,
		SourceStoragePlan: user.StoragePlan,
	}
	preferenceJSON, errMarshal := json.Marshal(task)
	if tools.PrintError(errMarshal) {
		return false
	}
	resp, errNewRequest := http.Post(*args.TransporterUrl+"/task", "application/json", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return false
	}
	var response model.TaskResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return false
	}
	return true
}

func SyncFile(path string, userID string, oldPlan *model.StoragePlan, newPlan *model.StoragePlan) (*model.TaskResponse, bool) {
	task := model.TaskRequest{
		TaskType:               args.TaskTypeMigrate,
		UserID:                 userID,
		SourcePath:             path,
		SourceStoragePlan:      *oldPlan,
		DestinationPath:        path,
		DestinationStoragePlan: *newPlan,
	}
	preferenceJSON, errMarshal := json.Marshal(task)
	if tools.PrintError(errMarshal) {
		return nil, false
	}
	resp, errNewRequest := http.Post(*args.TransporterUrl+"/task", "application/json", bytes.NewReader(preferenceJSON))
	if tools.PrintError(errNewRequest) {
		return nil, false
	}
	var response model.TaskResponse
	errDecoder := json.NewDecoder(resp.Body).Decode(&response)
	if tools.PrintError(errDecoder) {
		return nil, false
	}
	return &response, true
}
