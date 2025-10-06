package user

import (
	"fmt"
	userDto "go-rest-setup/src/app-backoffice/user/dto"
	helper "go-rest-setup/src/lib/helpers"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *UserService
}

func NewController(userService *UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	users, err := c.userService.GetAll()
	if err != nil {
		return helper.Error(ctx, fiber.StatusInternalServerError, nil)
	}

	return helper.Success(ctx, users, fiber.StatusOK)
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var payload userDto.CreateUserDto

	if err := ctx.BodyParser(&payload); err != nil {
		return helper.Error(ctx, fiber.StatusUnprocessableEntity, map[string]string{"body": "invalid request"})
	}

	if err := helper.Validate.Struct(&payload); err != nil {
		return helper.Error(ctx, fiber.StatusUnprocessableEntity, helper.ParseValidationError(err))
	}

	createdUser, err := c.userService.Create(payload)
	if err != nil {
		if dupErr := helper.ParseDuplicateError(err); dupErr != nil {
			return helper.Error(ctx, fiber.StatusUnprocessableEntity, dupErr)
		}
		return helper.Error(ctx, fiber.StatusInternalServerError, nil)
	}

	return helper.Success(ctx, createdUser, fiber.StatusCreated)
}

func (c *UserController) SoftDelete(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return helper.Error(ctx, fiber.StatusNotFound, map[string]string{"id": fmt.Sprintf("invalid id %s", strconv.Itoa(id))})
	}

	if err := c.userService.SoftDelete(uint(id)); err != nil {
		return helper.Error(ctx, fiber.StatusNotFound, err)
	}

	return helper.SuccessVoid(ctx, fiber.StatusNoContent)

}
