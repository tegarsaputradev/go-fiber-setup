package appbackoffice

import (
	"go-rest-setup/src/app-backoffice/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BackofficeContainer struct {
	RedisClient    *redis.Client
	UserController *user.UserController
}

func NewBackofficeContainer(db *gorm.DB, redis *redis.Client) *BackofficeContainer {
	return &BackofficeContainer{
		RedisClient:    redis,
		UserController: user.NewController(user.NewService(db)),
	}
}
