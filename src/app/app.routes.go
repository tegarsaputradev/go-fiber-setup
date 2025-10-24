package app

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, c *AppController) {
	r.Get("/", c.GetHello)
}
