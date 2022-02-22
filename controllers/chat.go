package controllers

import (
	"context"
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
	return redisClient
}

func initSession(session *melody.Session) string {
	id := uuid.New().String()
	session.Set(chatKey, id)
	return id
}

func getSessionID(session *melody.Session) string {
	if id, isExist := session.Get(chatKey); isExist {
		return id.(string)
	}
	return initSession(session)
}

func addToWaitList(id string, redisClient *redis.Client) error {
	return redisClient.LPush(context.Background(), chatWait, id).Err()
}

func getWaitFirstKey(redisClient *redis.Client) (string, error) {
	return redisClient.LPop(context.Background(), chatWait).Result()
}

func createChat(id1, id2 string, redisClient *redis.Client) {
	redisClient.Set(context.Background(), id1, id2, 0)
	redisClient.Set(context.Background(), id2, id1, 0)
}

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

		msgJ := NewMessageJ(string(msg), GetUsername(session))
		msgController.InsertMessage(msgJ)

		msg = NewMessage(string(msg), GetUsername(session))
		//socket chat filtered with sessionID and only allow 1 on 1 talk
		melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(chatKey)
			return compareID == chatTo || compareID == id
		})
	})
}

//when user connect the chat
func (chatController ChatController) HandleConnect(socketTool *mdDB.RedisTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient

	melodyNew.HandleConnect(func(session *melody.Session) {
		id := initSession(session)
		//check if 2 people are waiting to chat
		if key, err := getWaitFirstKey(redisClient); err == nil && key != "" {
			//create 1 on 1 chat
			createChat(id, key, redisClient)
			msg := NewMessage("Match success", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(chatKey)
				return compareID == id || compareID == key
			})
		} else {
			// if not then add to the queue and wait
			addToWaitList(id, redisClient)
			msg := NewMessage("Wait for matching...", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(chatKey)
				return compareID == id
			})
		}
	})
}

//when user disconnect
func (chatController ChatController) HandleClose(socketTool *mdDB.RedisTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient

	melodyNew.HandleClose(func(session *melody.Session, i int, s string) error {
		id := getSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg := NewMessage(GetUsername(session)+" left the room", "Server")
		removeChat(id, chatTo, redisClient)
		return melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(chatKey)
			return compareID == chatTo
		})
	})
}
