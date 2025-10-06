package appbackoffice

import "go-rest-setup/src/app-backoffice/user"

type BackofficeContainer struct {
	UserController *user.UserController
}

func NewBackofficeContainer() *BackofficeContainer {
	return &BackofficeContainer{
		UserController: user.NewController(user.NewService()),
	}
}
