package database

import "go.mongodb.org/mongo-driver/mongo"

type MongoTool struct {
	MongoClient *mongo.Client
	Database    *mongo.Database
	CollName    *mongo.Collection
}

type MongoConfig struct {
	Url        string
	Password   string
	Database   string
	Collection string
	Port       int
}
