package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/harshgupta9473/webhookDelivery/config"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

// init reddis
func InitRedis() {
	c := config.LoadConfig()

	// a new Redis client instance
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort),
		Password: c.RedisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	log.Println("Redis connected.")
}
