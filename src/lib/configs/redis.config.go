package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	c := EnvModule()

	db, err := strconv.Atoi(c.Redis.Database)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = client.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis (%s): %v", client.Options().Addr, err)
	}

	log.Printf("✅ Connected to Redis at %s [DB: %d]", client.Options().Addr, db)
	return client
}
