package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CustomResponse(ctx *fiber.Ctx) error {
	if err := ctx.Next(); err != nil {
		return err
	}

	if len(ctx.Response().Body()) > 0 {
		return nil
	}

	resp := ctx.Locals("response")
	if resp == nil {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"status_code": http.StatusOK,
			"message":     http.StatusText(http.StatusOK),
			"data":        nil,
		})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
