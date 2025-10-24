package authDto

type LoginUsernameDto struct {
	Username string `json:"username" validate:"required" example:"superadmin"`
	Password string `json:"password" validate:"required" example:"ggwp"`
}
