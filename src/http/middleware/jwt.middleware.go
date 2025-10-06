package middleware

import (
	"context"
	"fmt"
	config "go-rest-setup/src/lib/configs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

func JwtProtected(redis *redis.Client) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(config.EnvModule().JWT.Secret),
		ContextKey:    "user",
		SigningMethod: "HS256",
		ErrorHandler:  jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			ctx := context.Background()
			authHeader := c.Get("Authorization")
			if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  "error",
					"message": "missing or invalid token",
				})
			}
			userToken := authHeader[7:]

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userID := fmt.Sprintf("%v", claims["user_id"])

			sval, err := redis.Get(ctx, fmt.Sprintf("AUTH:%s", userID)).Result()
			if err != nil || sval != userToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status_code": fiber.StatusUnauthorized,
					"message":     "session expired or logged out",
				})
			}

			return c.Next()
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status_code": fiber.StatusUnauthorized,
		"message":     http.StatusText(fiber.StatusUnauthorized),
	})
}
