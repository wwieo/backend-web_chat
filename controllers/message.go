package controllers

import (
	mdDB "backend-web_chat/model/database"
	mdMsg "backend-web_chat/model/message"
	"fmt"
	"log"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	utilsController = NewUtilsController()
	mongoConfig     = *new(mdDB.MongoConfig)
)

func init() {
	config := utilsController.GetConfig()
	mongoConfig = mdDB.MongoConfig{
		Url:        config.GetString("mongo.url"),
		Password:   config.GetString("mongo.password"),
		Database:   config.GetString("mongo.database"),
		Collection: config.GetString("mongo.collection"),
		Port:       config.GetInt("mongo.port"),
	}
}

type MessageController struct {
	mongoTool *mdDB.MongoTool
}

func NewMessageController(mongoTool *mdDB.MongoTool) *MessageController {
	return &MessageController{mongoTool}
}

func SetMongoClient() *mongo.Client {
	URI := fmt.Sprintf("%s:%d", mongoConfig.Url, mongoConfig.Port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB successfully connected and pinged.")

	return client
}

func (msgController MessageController) InsertMessage(message *mdMsg.Message) {
	msgController.mongoTool.CollName.InsertOne(context.Background(), message)
}

func (msgController MessageController) NewMessage(body string, sender string) (msg *mdMsg.Message) {
	msg = &mdMsg.Message{
		IP:     utilsController.GetUserIP(),
		ID:     utilsController.GetMsgID(),
		Body:   body,
		Sender: sender,
		Time:   time.Now(),
	}
	return msg
}
