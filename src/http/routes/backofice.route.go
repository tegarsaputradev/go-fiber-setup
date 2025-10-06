package routes

import (
	appbackoffice "go-rest-setup/src/app-backoffice"
	"go-rest-setup/src/app-backoffice/user"
	"go-rest-setup/src/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterBackofficeRoutes(app *fiber.App, container *appbackoffice.BackofficeContainer) {
	backoffice := app.Group(`/api/v1/backoffice`, middleware.JwtProtected(container.RedisClient))
	user.RegisterRoutes(backoffice, container.UserController)
}
