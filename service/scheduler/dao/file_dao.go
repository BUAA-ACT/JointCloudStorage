package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shaoliyin.me/jcspan/entity"
	"shaoliyin.me/jcspan/tools"
)

func GetFile(col *mongo.Collection, fid string) (entity.File, error) {
	//col := d.client.Database(d.database).Collection(d.fileCollection)
	var file entity.File
	if colErr := VerifyCollection(col); colErr != nil {
		return file, colErr
	}
	err := col.FindOne(context.TODO(), bson.M{"file_id": fid}).Decode(&file)
	return file, err
}

func InsertFiles(col *mongo.Collection, files []entity.File) error {
	if colErr := VerifyCollection(col); colErr != nil {
		return colErr
	}
	fs := make([]interface{}, len(files))
	for i := range files {
		fs[i] = files[i]
	}
	//col := d.client.Database(d.database).Collection(d.fileCollection)
	for _, file := range files {
		_, err := col.UpdateOne(
			context.TODO(),
			bson.M{
				"file_id": file.FileID,
			},
			bson.D{
				{"$set", file},
			},
			&options.UpdateOptions{
				Upsert: tools.Bool2Pointer(true),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteFiles(col *mongo.Collection, files []entity.File) error {
	if colErr := VerifyCollection(col); colErr != nil {
		return colErr
	}
	var fs []string
	for _, v := range files {
		fs = append(fs, v.FileID)
	}
	//col := d.client.Database(d.database).Collection(d.fileCollection)
	_, err := col.DeleteMany(
		context.TODO(),
		bson.M{
			"file_id": bson.M{"$in": fs},
		},
	)
	if err != nil {
		return err
	}

	return nil
}