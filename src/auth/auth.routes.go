package auth

import (
	"go-rest-setup/src/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r fiber.Router, c *AuthController, redisClient *redis.Client) {
	auth := r.Group(`/auth`)
	auth.Post(`/register`, c.Register)
	auth.Post(`/login`, c.Login)
	auth.Post(`/logout/:id`, c.Logout)
	auth.Get(`/get-me`, middleware.JwtProtected(redisClient), c.GetMe)

}
