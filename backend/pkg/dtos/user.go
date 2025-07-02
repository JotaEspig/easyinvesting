package dtos

import (
	"easyinvesting/pkg/utils"

	"github.com/microcosm-cc/bluemonday"
)

type UserDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Omit password in JSON responses
}

func (u *UserDTO) Sanitize(policy *bluemonday.Policy) {
	u.Email = policy.Sanitize(u.Email)
}

func (u UserDTO) IsValid() bool {
	return u.Email != "" && u.Password != "" && utils.IsValidEmail(u.Email)
}

func (u UserDTO) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
	}
}
