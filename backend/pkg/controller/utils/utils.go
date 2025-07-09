package utils

import (
	"easyinvesting/pkg/types"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromContext(c echo.Context) (*types.JWTClaims, error) {
	userAuth := c.Get("user")
	if userAuth == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	token, ok := userAuth.(*jwt.Token)
	if !ok || token == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok || claims == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	return claims, nil
}
