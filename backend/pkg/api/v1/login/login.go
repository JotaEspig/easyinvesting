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
)

func login(c echo.Context) error {
	u := user.User{}
	err := json.NewDecoder(c.Request().Body).Decode(&u)
	if !u.IsValid() || err != nil {
		return c.JSON(http.StatusBadRequest, types.JsonMap{
			"message": "some user field may be missing or invalid",
		})
	}
	tx, err := config.DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.JsonMap{
			"message": "creating of database transaction failed. Try again",
		})
	}
	defer tx.Rollback()
	savedUser, err := user.FindUserByUsername(tx, u.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{
			"message": "unauthorized",
		})
	}
	if !savedUser.ValidateLogin(u.Username, u.Password) {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{
			"message": "unauthorized",
		})
	}
	claims := &types.JWTClaims{
		UserID: savedUser.ID,
		Type:   savedUser.Type,
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
