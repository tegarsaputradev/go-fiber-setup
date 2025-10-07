package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r fiber.Router, c *AuthController) {
	auth := r.Group(`/auth`)
	auth.Post(`/login`, c.Login)
	auth.Post(`/logout/:id`, c.Logout)
	auth.Get(`/get-me/:id`, c.GetMe)
}
