package controllers

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	mdDB "backend-web_chat/model/database"
	utils "backend-web_chat/utils"
)

const (
	KEY  = "chat_id"
	WAIT = "wait"
)

func SetRedisClient() *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Password: "", no password set
		DB: 0, // use default DB
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis error:", err)
	}
	return RedisClient
}

func InitSession(session *melody.Session) string {
	id := uuid.New().String()
	session.Set(KEY, id)
	return id
}

func GetSessionID(session *melody.Session) string {
	if id, isExist := session.Get(KEY); isExist {
		return id.(string)
	}
	return InitSession(session)
}

func AddToWaitList(id string, redisClient *redis.Client) error {
	return redisClient.LPush(context.Background(), WAIT, id).Err()
}

func GetWaitFirstKey(redisClient *redis.Client) (string, error) {
	return redisClient.LPop(context.Background(), WAIT).Result()
}

func CreateChat(id1, id2 string, redisClient *redis.Client) {
	redisClient.Set(context.Background(), id1, id2, 0)
	redisClient.Set(context.Background(), id2, id1, 0)
}

func RemoveChat(id1, id2 string, redisClient *redis.Client) {
	redisClient.Del(context.Background(), id1, id2)
}

//when user request the route
func HandleRequest(socketTool *mdDB.RedisTool) func(c *gin.Context) {
	return func(c *gin.Context) {
		socketTool.Melody.HandleRequest(c.Writer, c.Request)
	}
}

//when user send the message
func HandleMessage(socketTool *mdDB.RedisTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient

	melodyNew.HandleMessage(func(session *melody.Session, msg []byte) {
		id := GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()

		// omsg := OriginMessage(string(msg), utils.GetUsername(session))
		// controller := controllers.NewMessageContoller(mongoTool)
		// controller.InsertMessage(omsg)

		msg = NewMessage(string(msg), utils.GetUsername(session))
		//socket chat filtered with sessionID and only allow 1 on 1 talk
		melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo || compareID == id
		})
	})
}

//when user connect the chat
func HandleConnect(socketTool *mdDB.RedisTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient

	melodyNew.HandleConnect(func(session *melody.Session) {
		id := InitSession(session)
		//check if 2 people are waiting to chat
		if key, err := GetWaitFirstKey(redisClient); err == nil && key != "" {
			//create 1 on 1 chat
			CreateChat(id, key, redisClient)
			msg := NewMessage("Match success", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(KEY)
				return compareID == id || compareID == key
			})
		} else {
			// if not then add to the queue and wait
			AddToWaitList(id, redisClient)
			msg := NewMessage("Wait for matching...", "Server")
			melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(KEY)
				return compareID == id
			})
		}
	})
}

//when user disconnect
func HandleClose(socketTool *mdDB.RedisTool) {
	melodyNew := socketTool.Melody
	redisClient := socketTool.RedisClient

	melodyNew.HandleClose(func(session *melody.Session, i int, s string) error {
		id := GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg := NewMessage(utils.GetUsername(session)+" left the room", "Server")
		RemoveChat(id, chatTo, redisClient)
		return melodyNew.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo
		})
	})
}
