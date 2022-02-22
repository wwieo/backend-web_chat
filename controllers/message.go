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

type MessageController struct {
	mongoTool *mdDB.MongoTool
}

func NewMessageContoller(mongoTool *mdDB.MongoTool) *MessageController {
	return &MessageController{mongoTool}
}

func SetMongoClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
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

	//ISODate
	localTime, _ := time.LoadLocation("Local")
	msg = &mdMsg.Message{
		IP:     GetUserIP(),
		ID:     GetMsgID(),
		Body:   body,
		Sender: sender,
		Time:   time.Now().In(localTime),
	}
	return msg
}
