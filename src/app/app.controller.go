package app

import (
	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	appService *AppService
}

func NewController(appService *AppService) *AppController {
	return &AppController{
		appService: appService,
	}
}

func (c *AppController) GetHello(ctx *fiber.Ctx) error {

	return ctx.SendString(c.appService.Hello())

}
