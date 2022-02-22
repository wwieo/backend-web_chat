package controllers

import (
	mdDB "backend-web_chat/model/database"
	mdMsg "backend-web_chat/model/message"
	"backend-web_chat/utils"
	"encoding/json"

	"fmt"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageController struct {
	mongoTool *mdDB.MongoTool
}

func NewMessageContoller(mongoTool *mdDB.MongoTool) *MessageController {
	return &MessageController{mongoTool}
}

func (msgController MessageController) InsertMessage(message *mdMsg.Message) {
	fmt.Print(message)
	msgController.mongoTool.CollName.InsertOne(context.Background(), message)
}

func SetMongoClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func NewMessage(body string, sender string) []byte {
	result, _ := json.Marshal(
		&mdMsg.Message{
			IP:     utils.GetUserIP(),
			ID:     utils.GetMsgID(),
			Body:   body,
			Sender: sender,
		})
	return result
}
