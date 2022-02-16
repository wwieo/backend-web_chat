package match

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
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
