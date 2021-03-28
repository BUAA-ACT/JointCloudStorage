package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckClient(client *mongo.Client,connectionOptions *options.ClientOptions)(error){
	 err:=client.Ping(context.TODO(),nil)
	 if err!=nil{
	 	client,err=mongo.Connect(context.TODO(),connectionOptions)
	 	if err!=nil{
	 		return err
		}
	 }
	 return nil
}