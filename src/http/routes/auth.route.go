package routes

import (
	"go-rest-setup/src/auth"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, c *auth.AuthController) {
	v1 := app.Group(`/api/v1`)
	auth.RegisterRoutes(v1, c)
}
