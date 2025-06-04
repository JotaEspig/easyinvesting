package user

import (
	"easyinvesting/pkg/types"
	"easyinvesting/pkg/utils"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"username"`
	hashedPassword string
}

func (u *User) SetPassword(password string) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	utils.HandleErr(err, "Failed to hash password")
	u.hashedPassword = string(hashed)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.hashedPassword), []byte(password))
	return err == nil // If the password matches, err will be nil
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.hashedPassword != "" && utils.IsValidEmail(u.Email)
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
