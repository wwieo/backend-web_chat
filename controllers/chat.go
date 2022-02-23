package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	mdDB "backend-web_chat/model/database"
)

const (
	chatKey  = "chat_id"
	chatWait = "wait"
)

type ChatController struct {
	socketTool *mdDB.RedisTool
}

func NewChatContoller(socketTool *mdDB.RedisTool) *ChatController {
	return &ChatController{socketTool}
}

func SetRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Password: "", no password set
		DB: 0, // use default DB
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis error:", err)
	}
	fmt.Println("Redis successfully connected and pinged.")
	return redisClient
}

//create a new session id for chatting
func initSession(session *melody.Session) string {
	id := uuid.New().String()
	session.Set(chatKey, id)
	return id
}

//when user connect the server, then return a session id
func getSessionID(session *melody.Session) string {
	if id, isExist := session.Get(chatKey); isExist {
		return id.(string)
	}
	return initSession(session)
}

//add people with id to wait list
func addToWaitList(id string, redisClient *redis.Client) error {
	return redisClient.LPush(context.Background(), chatWait, id).Err()
}

//pop a wating people's id
func getWaitFirstKey(redisClient *redis.Client) (string, error) {
	return redisClient.LPop(context.Background(), chatWait).Result()
}

//when 2 people pair successully, create a one-to-one chat
func createChat(id1, id2 string, redisClient *redis.Client) {
	redisClient.Set(context.Background(), id1, id2, 0)
	redisClient.Set(context.Background(), id2, id1, 0)
}

//chat ended
func removeChat(id1, id2 string, redisClient *redis.Client) {
	redisClient.Del(context.Background(), id1, id2)
}

//when user request the route
func HandleRequest(socketTool *mdDB.RedisTool) func(c *gin.Context) {
	return func(c *gin.Context) {
		socketTool.Melody.HandleRequest(c.Writer, c.Request)
	}
}

//when user send the message
func (chatController ChatController) HandleMessage(socketTool *mdDB.RedisTool, mongoTool *mdDB.MongoTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient
	msgController := NewMessageContoller(mongoTool)

	melodyNew.HandleMessage(func(session *melody.Session, msg []byte) {
		id := getSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()

		newMsg := msgController.NewMessage(string(msg), GetUsername(session))

		msgController.InsertMessage(newMsg)
		msg, _ = json.Marshal(newMsg)

		//socket chat filtered by sessionID and only allow one-to-one talk
		melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(chatKey)
			return compareID == chatTo || compareID == id
		})
	})
}

//when user connect the chat
func (chatController ChatController) HandleConnect(socketTool *mdDB.RedisTool, mongoTool *mdDB.MongoTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient
	msgController := NewMessageContoller(mongoTool)

	melodyNew.HandleConnect(func(session *melody.Session) {
		id := initSession(session)
		//get people in waitlist
		if key, err := getWaitFirstKey(redisClient); err == nil && key != "" {
			createChat(id, key, redisClient)

			newMsg := msgController.NewMessage("Match success", "Server")
			msgController.InsertMessage(newMsg)
			msg, _ := json.Marshal(newMsg)
			fmt.Println(newMsg.Time)

			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(chatKey)
				return compareID == id || compareID == key
			})
		} else {
			//if not then add to the queue and wait
			addToWaitList(id, redisClient)
			msg, _ := json.Marshal(msgController.NewMessage("Wait for matching...", "Server"))

			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(chatKey)
				return compareID == id
			})
		}
	})
}

//when user disconnect
func (chatController ChatController) HandleClose(socketTool *mdDB.RedisTool, mongoTool *mdDB.MongoTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient
	msgController := NewMessageContoller(mongoTool)

	melodyNew.HandleClose(func(session *melody.Session, i int, s string) error {
		id := getSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()

		newMsg := msgController.NewMessage(GetUsername(session)+" left the room", "Server")
		msgController.InsertMessage(newMsg)
		msg, _ := json.Marshal(newMsg)

		removeChat(id, chatTo, redisClient)
		return melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(chatKey)
			return compareID == chatTo
		})
	})
}
