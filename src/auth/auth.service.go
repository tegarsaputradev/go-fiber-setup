package auth

import (
	"context"
	"errors"
	"fmt"
	authDto "go-rest-setup/src/auth/dto"
	"go-rest-setup/src/database/models"
	config "go-rest-setup/src/lib/configs"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type JwtClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
}

func NewService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(payload authDto.LoginUsernameDto) (string, *models.User, error) {
	var user *models.User
	if err := config.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		return "", nil, fmt.Errorf("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	claims := JwtClaims{
		UserID:   user.ID,
		Username: payload.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.EnvModule().JWT.Secret))
	if err != nil {
		return "", nil, err
	}

	ctx := context.Background()
	key := fmt.Sprintf("AUTH:%d", user.ID)
	value := token
	expiration := 24 * time.Hour

	if err := config.RedisClient.Set(ctx, key, value, expiration).Err(); err != nil {
		return "", nil, fmt.Errorf("failed to store token in Redis: %v", err)
	}

	return token, user, nil

}

func (s *AuthService) GetMe(id uint) (*models.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf(`AUTH:%d`, id)

	userSession, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil || userSession == "" {
		return nil, fmt.Errorf("session not found or already loged out")
	}

	var user models.User

	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("session is anonymous")
	}

	return &user, nil

}

func (s *AuthService) Logout(id uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("AUTH:%d", id)

	if err := config.RedisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}
