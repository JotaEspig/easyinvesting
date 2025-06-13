package user

import (
	"easyinvesting/pkg/types"
	"easyinvesting/pkg/utils"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `json:"email" gorm:"uniqueIndex;not null"`
	Password       string `json:"password,omitempty" gorm:"-"` // Omitido na serialização JSON
	HashedPassword string `json:"-" gorm:"not null"`           // Campo para armazenar a senha hasheada, não exposto na API
}

func (u *User) SetPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	utils.HandleErr(err, "Failed to hash password")
	u.HashedPassword = string(hashed)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil // If the password matches, err will be nil
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.Password != "" && utils.IsValidEmail(u.Email)
}

func (u *User) Sanitize(policy *bluemonday.Policy) {
	u.Email = policy.Sanitize(u.Email)
}

func (u *User) ToMap() types.JsonMap {
	m := make(types.JsonMap)
	m["id"] = u.ID
	m["email"] = u.Email
	return m
}
