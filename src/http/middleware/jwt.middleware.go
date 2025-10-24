package middleware

import (
	"context"
	"fmt"
	config "go-rest-setup/src/lib/configs"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type UserSession struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ErrorData struct {
	Message string `json:"message"`
}

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
				return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
					StatusCode: fiber.StatusUnauthorized,
					Message:    http.StatusText(http.StatusUnauthorized),
					Data: ErrorData{
						Message: "Session expired or logged out.",
					},
				})
			}

			return c.Next()
		},
	})
}

func GetSessionUser(ctx *fiber.Ctx) (*UserSession, error) {
	userSession := ctx.Locals("user")

	if userSession == nil {
		return nil, fmt.Errorf("no user session found in context")
	}

	token, ok := userSession.(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("invalid token type in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT claims")
	}

	userIDStr := fmt.Sprintf("%v", claims["user_id"])
	username := fmt.Sprintf("%v", claims["username"])
	email := fmt.Sprintf("%v", claims["email"])

	idInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %v", err)
	}

	session := &UserSession{
		UserID:   uint(idInt),
		Username: username,
		Email:    email,
	}

	return session, nil
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status_code": fiber.StatusUnauthorized,
		"message":     http.StatusText(fiber.StatusUnauthorized),
	})
}
