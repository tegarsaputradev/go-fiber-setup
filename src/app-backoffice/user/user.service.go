package user

import (
	userDto "go-rest-setup/src/app-backoffice/user/dto"
	"go-rest-setup/src/database/models"
	config "go-rest-setup/src/lib/configs"

	"gorm.io/gorm"
)

type UserService struct{}

func NewService() *UserService {
	return &UserService{}
}

func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) Create(payload userDto.CreateUserDto) (*models.User, error) {
	user := &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindOne(condition map[string]interface{}, relations []string) (*models.User, error) {
	var user models.User

	query := config.DB.Model(&models.User{})
	for _, v := range relations {
		query = query.Preload(v)
	}

	if err := query.Where(condition).First(user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) FindOneById(id uint) (*models.User, error) {
	return s.FindOne(map[string]interface{}{"id": id}, nil)
}

func (s *UserService) SoftDelete(id uint) error {
	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
