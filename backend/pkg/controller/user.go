package controller

import (
	"easyinvesting/config"
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/service"
	"easyinvesting/pkg/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		userDTO := new(dto.UserDTO)
		if err := c.Bind(userDTO); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		user, err := controller.userService.Login(userDTO.Email, userDTO.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}

		claims := &types.JWTClaims{
			UserID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(config.SecretKey()))
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Failed to generate token"},
			)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token": signedToken,
			"user":  user.ToMap(),
		})
	}
}

func (controller UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		userDTO := new(dto.UserDTO)
		if err := c.Bind(userDTO); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		userDTO.Sanitize(config.StrictPolicy)
		if !userDTO.IsValid() {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user data"})
		}

		if existingUser, _ := controller.userService.Login(userDTO.Email, ""); existingUser != nil {
			return c.JSON(http.StatusConflict, map[string]string{"message": "User already exists"})
		}

		if err := controller.userService.Save(userDTO); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusCreated, userDTO)
	}
}
