package userDto

type CreateUserDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
	Username string `json:"username" validate:"required"`
}
