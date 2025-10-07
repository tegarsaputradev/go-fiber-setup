package file

import (
	"context"
	helper "go-rest-setup/src/lib/helpers"

	"github.com/gofiber/fiber/v2"
)

type FileController struct {
	fileService *FileService
}

func NewController(fs *FileService) *FileController {
	return &FileController{
		fileService: fs,
	}
}

func (s *FileController) Upload(ctx *fiber.Ctx) error {
	filHeader, err := ctx.FormFile("file")
	if err != nil {
		return helper.Error(ctx, fiber.StatusBadRequest, map[string]string{
			"file": "file not found in form-data",
		})
	}

	file, err := s.fileService.Upload(context.Background(), filHeader, "")

	if err != nil {
		return helper.Error(ctx, fiber.StatusInternalServerError, map[string]string{
			"upload": err.Error(),
		})
	}

	return helper.Success(ctx, file, fiber.StatusCreated)

}
