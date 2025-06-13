package other

import (
	"easyinvesting/pkg/types"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, types.JsonMap{"message": "Hello from Go Echo!"})
}

func helloAuth(c echo.Context) error {
	userAuth := c.Get("user")
	if userAuth == nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	var token *jwt.Token
	var ok bool
	var claims *types.JWTClaims
	if token, ok = userAuth.(*jwt.Token); !ok || token == nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}
	if claims, ok = token.Claims.(*types.JWTClaims); !ok || claims == nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	user := types.JsonMap{
		"id": claims.UserID,
	}
	return c.JSON(http.StatusOK, types.JsonMap{"message": "Hello from Go Echo with Auth!", "user": user})
}
