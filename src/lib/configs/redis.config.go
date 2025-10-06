package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	db, err := strconv.Atoi(EnvModule().Redis.Database)
	if err != nil {
		db = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", EnvModule().Redis.Host, EnvModule().Redis.Port),
		Password: EnvModule().Redis.Password,
		DB:       db,
	})

	ctx := context.Background()
	if _, err = RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	log.Println("✅ Connected to Redis")
}
