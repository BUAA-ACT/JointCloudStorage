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

	a.Databases["aaa"].Collections["bbb"] = nil
	b:= ClientConfig{
		Client:    nil,
		Databases: nil,
	}
	for databaseName, database := range a.Databases {
		fmt.Println(databaseName)
		fmt.Println(database)
	}
	fmt.Println()
	for databaseName, database := range b.Databases {
		fmt.Println(databaseName)
		fmt.Println(database)
	}

	ttt := config{aaa: 1}
	pc := Config{config: &ttt}

	os.Exit(1)
}
