package models

import (
	"easyinvesting/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex;not null"`
	HashedPassword string `gorm:"not null"`
}

func (u *User) SetPassword(password string) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	utils.HandleErr(err, "Failed to hash password")
	u.HashedPassword = string(hashed)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.HashedPassword != "" && utils.IsValidEmail(u.Email)
}
