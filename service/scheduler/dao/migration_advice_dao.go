package dao

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shaoliyin.me/jcspan/entity"
)

func InsertMigrationAdvice(col *mongo.Collection, adv entity.MigrationAdvice) error {
	//col := d.client.Database(d.database).Collection(d.migrationAdvice)
	if colErr := VerifyCollection(col); colErr != nil {
		return colErr
	}
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