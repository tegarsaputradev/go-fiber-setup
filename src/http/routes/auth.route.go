package routes

import (
	"go-rest-setup/src/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterAuthRoutes(app *fiber.App, c *auth.AuthController, redis *redis.Client) {
	v1 := app.Group(`/api/v1`)
	auth.RegisterRoutes(v1, c, redis)
}
