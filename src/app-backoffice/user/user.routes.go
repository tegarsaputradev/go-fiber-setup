package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, c *UserController) {
	users := r.Group("/users")
	users.Get("/", c.GetAll)
	users.Post("/", c.Create)
	users.Delete("/:id", c.SoftDelete)
}
