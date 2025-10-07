package file

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, c *FileController) {
	file := router.Group(`/file`)
	file.Post(`/upload`, c.Upload)
}
