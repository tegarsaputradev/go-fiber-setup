package auth

import (
	authDto "go-rest-setup/src/auth/dto"
	"go-rest-setup/src/database/models"
	helper "go-rest-setup/src/lib/helpers"

	"github.com/gofiber/fiber/v2"
)

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type AuthController struct {
	authService *AuthService
}

func NewController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (s *AuthController) Login(ctx *fiber.Ctx) error {
	var payload authDto.LoginUsernameDto

	if err := ctx.BodyParser(&payload); err != nil {
		return helper.Error(ctx, fiber.StatusUnprocessableEntity, map[string]string{"body": "invalid request"})
	}

	if err := helper.Validate.Struct(&payload); err != nil {
		return helper.Error(ctx, fiber.StatusUnprocessableEntity, helper.ParseValidationError((err)))
	}

	token, user, err := s.authService.Login(payload)
	if err != nil {
		return helper.Error(ctx, fiber.StatusInternalServerError, nil)
	}

	responseData := LoginResponse{
		Token: token,
		User:  user,
	}

	return helper.Success(ctx, responseData, fiber.StatusAccepted)

}

func (s *AuthController) GetMe(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("id")
	if err != nil || userID <= 0 {
		return helper.Error(ctx, fiber.StatusBadRequest, map[string]string{"id": "invalid user ID"})
	}
	user, errGet := s.authService.GetMe(uint(userID))
	if errGet != nil {
		if errGet.Error() == "session not found or already loged out" {
			return helper.Error(ctx, fiber.StatusNotFound, map[string]string{"session": "session not found or already loged out"})
		}
		return helper.Error(ctx, fiber.StatusInternalServerError, errGet)
	}
	return helper.Success(ctx, user, fiber.StatusOK)
}

func (s *AuthController) Logout(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("id")
	if err != nil || userID <= 0 {
		return helper.Error(ctx, fiber.StatusBadRequest, map[string]string{
			"id": "invalid user ID",
		})
	}

	if err := s.authService.Logout(uint(userID)); err != nil {
		return helper.Error(ctx, fiber.StatusUnauthorized, map[string]string{
			"logout": err.Error(),
		})
	}

	return helper.Success(ctx, fiber.Map{
		"message": "Logged out successfully",
	}, fiber.StatusOK)
}
