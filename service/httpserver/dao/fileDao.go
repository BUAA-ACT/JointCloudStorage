package dao

import (
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *Dao) ListFiles(userId string, path string, isDir bool) (*[]model.File, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	var files []model.File = make([]model.File, 0)
	var filter interface{}
	// TODO time complex high !!!
	filterDir := bson.M{
		"owner": userId,
	}
	filterFile := bson.M{
		"owner":    userId,
		"filename": path,
	}
	if isDir {
		filter = filterDir
	} else {
		filter = filterFile
	}
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var file model.File
		err := result.Decode(&file)
		if tools.PrintError(err) {
			return nil, false
		}
		files = append(files, file)
	}
	return &files, true
}

func (d *Dao) CheckFileStatus(userId string, path string) (*[]model.File, bool) {
	col := d.client.Database(d.database).Collection(d.collection)
	filter := bson.M{
		"owner":    userId,
		"filename": path,
	}
	files := make([]model.File, 0)
	result, err := col.Find(context.TODO(), filter)
	if tools.PrintError(err) {
		return nil, false
	}
	for result.Next(context.TODO()) {
		var file model.File
		err := result.Decode(&file)
		if tools.PrintError(err) {
			return nil, false
		}
		files = append(files, file)
	}
	return &files, true
}
