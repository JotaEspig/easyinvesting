package login

import (
	"easyinvesting/config"
	"easyinvesting/pkg/models/user"
	"easyinvesting/pkg/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func login(c echo.Context) error {
	var u user.User

	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil || !u.IsValid() {
		if err != nil {
			c.Logger().Errorf("Failed to decode user: %v", err)
		} else {
			c.Logger().Errorf("Invalid user data: %v", u)
		}
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some user field may be missing or invalid",
		})
	}

	var savedUser user.User
	err := config.DB.Where("email = ?", u.Email).First(&savedUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Logger().Errorf("Login attempt for non-existent user: %s", u.Email)
		} else {
			c.Logger().Errorf("Database error during login: %v", err)
			return c.JSON(http.StatusInternalServerError, types.JsonMap{
				"message": "internal server error",
			})
		}

		return c.JSON(http.StatusUnauthorized, types.JsonMap{
			"message": "unauthorized",
		})
	}

	if !savedUser.CheckPassword(u.Password) {
		c.Logger().Errorf("Failed login attempt for user: %s", u.Email)
		return c.JSON(http.StatusUnauthorized, types.JsonMap{
			"message": "unauthorized",
		})
	}

	claims := &types.JWTClaims{
		UserID: savedUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.SecretKey()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.JsonMap{"token": t})
}

func signup(c echo.Context) error {
	var u user.User

	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil || !u.IsValid() {
		c.Logger().Errorf("Failed to decode user: %v", err)
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some user field may be missing or invalid",
		})
	}

	u.Sanitize(config.StrictPolicy)
	u.SetPassword()

	if err := config.DB.Create(&u).Error; err != nil {
		c.Logger().Errorf("Failed to create user: %v", err)
		if err == gorm.ErrDuplicatedKey {
			return c.JSON(http.StatusConflict, types.JsonMap{
				"message": "user already exists",
			})
		}

		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, u.ToMap())
}
