package controllers

import (
	mdDB "backend-web_chat/model/database"
	mdMsg "backend-web_chat/model/message"
	"encoding/json"
	"fmt"

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
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Mongo successfully connected and pinged.")

	return client
}

func (msgController MessageController) InsertMessage(message *mdMsg.Message) {
	msgController.mongoTool.CollName.InsertOne(context.Background(), message)
}

func NewMessage(body string, sender string) []byte {
	result, _ := json.Marshal(
		&mdMsg.Message{
			IP:     GetUserIP(),
			ID:     GetMsgID(),
			Body:   body,
			Sender: sender,
		})
	return result
}

func NewMessageJ(body string, sender string) *mdMsg.Message {
	return (&mdMsg.Message{
		IP:     GetUserIP(),
		ID:     GetMsgID(),
		Body:   body,
		Sender: sender,
	})
}
