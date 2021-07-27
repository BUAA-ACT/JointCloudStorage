package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) ListFiles(userID string, path string, isDir bool) (*[]model.File, bool) {
	var files []model.File = make([]model.File, 0)
	var filter interface{}
	filterDir := bson.M{
		"owner": userID,
		"file_name": bson.M{
			"$regex": path + "*",
		},
	}
	filterFile := bson.M{
		"owner":     userID,
		"file_name": path,
	}
	if isDir {
		filter = filterDir
	} else {
		filter = filterFile
	}
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var file model.File
		decodeErr := result.Decode(&file)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		files = append(files, file)
	}
	return &files, true
}

func (d *Dao) CheckFileStatus(userID string, path string) (*[]model.File, bool) {
	filter := bson.M{
		"file_id": userID + path,
	}
	files := make([]model.File, 0)
	result, findErr := d.collectionConnection.Find(context.TODO(), filter)
	if tools.PrintError(findErr) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var file model.File
		decodeErr := result.Decode(&file)
		if tools.PrintError(decodeErr) {
			return nil, false
		}
		files = append(files, file)
	}
	return &files, true
}
