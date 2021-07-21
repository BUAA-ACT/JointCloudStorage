package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) ListFiles(userID string, path string, isDir bool) (*[]model.File, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var files []model.File = make([]model.File, 0)
	var filter interface{}
	// TODO time complex high !!!
	filterDir := bson.M{
		"owner": userID,
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
	result, findErr := col.Find(context.TODO(), filter)
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
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"file_id": userID + path,
	}
	files := make([]model.File, 0)
	result, findErr := col.Find(context.TODO(), filter)
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
