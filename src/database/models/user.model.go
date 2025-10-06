package models

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel

	Name             string `json:"name" gorm:"not null"`
	Username         string `json:"username" gorm:"not null;unique"`
	Email            string `json:"email" gorm:"not null;unique"`
	Password         string `json:"-" gorm:"not null"`
	previousPassword string `gorm:"-"`
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.previousPassword = u.Password
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password != "" && u.Password != u.previousPassword {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		Password string `json:"password,omitempty"`
	}{
		Alias: (*Alias)(u),
	})
}
