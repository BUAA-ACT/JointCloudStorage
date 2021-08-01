package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)


type CollectionConfig struct {
	CollectionHandler *mongo.Collection
}

type DatabaseConfig struct {
	DatabaseHandler *mongo.Database
	Collections map[string]CollectionConfig
}

type ClientConfig struct {
	Client *mongo.Client
	Databases map[string]DatabaseConfig
}

type config struct {
	aaa int
}

type Config struct {
	config *config
}


func test() {
	a := ClientConfig{
		Client:    nil,
		Databases: map[string]DatabaseConfig{},
	}
	a.Databases["aaa"] = DatabaseConfig{
		DatabaseHandler: nil,
		Collections: map[string]CollectionConfig{},
	}
	kk1 := a.Databases
	fmt.Println(kk1)
	hh1 := kk1["aaa"].Collections
	fmt.Println(hh1)
	cc := a.Databases["kjhdskfjhsd"]

	fmt.Println(a.Databases["sdhajksd"])
	fmt.Println(cc)
	fmt.Println("---------------")
	b:= ClientConfig{
		Client:    nil,
		Databases: nil,
	}
	for databaseName, database := range a.Databases {
		fmt.Println(databaseName)
		fmt.Println(database)
	}
	kk := b.Databases
	fmt.Println(kk)
	fmt.Println(kk["aaa"])
	kk["test"] = DatabaseConfig{
		DatabaseHandler: nil,
		Collections: map[string]CollectionConfig{},
	}
	for databaseName, database := range b.Databases {
		fmt.Println(databaseName)
		fmt.Println(database)
	}

	ttt := config{aaa: 1}
	pc := Config{config: &ttt}
	fmt.Println(pc)
	os.Exit(1)
}
