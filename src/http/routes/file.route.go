package routes

import (
	"go-rest-setup/src/core/file"

	"github.com/gofiber/fiber/v2"
)

func RegisterFileRoute(app *fiber.App, c *file.FileController) {
	v1 := app.Group(`/api/v1`)
	file.RegisterRoutes(v1, c)
}
