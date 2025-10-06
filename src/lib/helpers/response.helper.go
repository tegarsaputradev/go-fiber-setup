package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

func SuccessVoid(ctx *fiber.Ctx, code int) error {
	return ctx.Status(code).JSON(Response{
		StatusCode: code,
		Message:    http.StatusText(code),
	})
}

func Success(ctx *fiber.Ctx, data interface{}, code int) error {
	return ctx.Status(code).JSON(Response{
		StatusCode: code,
		Message:    http.StatusText(code),
		Data:       data,
	})
}

func Error(ctx *fiber.Ctx, code int, errs interface{}) error {
	return ctx.Status(code).JSON(Response{
		StatusCode: code,
		Message:    http.StatusText(code),
		Errors:     errs,
	})
}
