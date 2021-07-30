package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Dao) GetFile(fid string) (File, error) {
	col := d.client.Database(d.database).Collection(d.fileCollection)

	var file File
	err := col.FindOne(context.TODO(), bson.M{"file_id": fid}).Decode(&file)
	return file, err
}

func (d *Dao) InsertFiles(files []File) error {
	fs := make([]interface{}, len(files))
	for i := range files {
		fs[i] = files[i]
	}

	col := d.client.Database(d.database).Collection(d.fileCollection)
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
				Upsert: bool2pointer(true),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dao) DeleteFiles(files []File) error {
	var fs []string
	for _, v := range files {
		fs = append(fs, v.FileID)
	}

	col := d.client.Database(d.database).Collection(d.fileCollection)
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
