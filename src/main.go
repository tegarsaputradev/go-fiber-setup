package main

import (
	appbackoffice "go-rest-setup/src/app-backoffice"
	"go-rest-setup/src/auth"
	"go-rest-setup/src/http/routes"
	config "go-rest-setup/src/lib/configs"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitDatabase()
	config.InitRedis()

	app := fiber.New()

	routes.RegisterBackofficeRoutes(app, appbackoffice.NewBackofficeContainer())
	routes.RegisterAuthRoutes(app, auth.NewController(auth.NewService()))

	if err := app.Listen(":" + config.EnvModule().Server.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
