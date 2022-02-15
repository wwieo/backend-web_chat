package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func Run() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Password: "", no password set
		DB: 0, // use default DB
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err == nil {
		log.Println(pong)
	} else {
		log.Fatal("no:", err)
	}
}
