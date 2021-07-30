package dao

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Dao) InsertMigrationAdvice(adv MigrationAdvice) error {
	col := d.client.Database(d.database).Collection(d.migrationAdvice)
	_, err := col.ReplaceOne(
		context.TODO(),
		bson.M{
			"user_id": adv.UserId,
		},
		adv,
		&options.ReplaceOptions{Upsert: aws.Bool(true)},
	)
	if err != nil {
		return err
	}

	return nil
}
