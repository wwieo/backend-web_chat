package database

import (
	"github.com/go-redis/redis/v8"
	"gopkg.in/olahol/melody.v1"
)

type RedisTool struct {
	Melody      *melody.Melody
	RedisClient *redis.Client
}

type RedisConfig struct {
	Url      string
	Password string
	Port     int
	DB       int
}
