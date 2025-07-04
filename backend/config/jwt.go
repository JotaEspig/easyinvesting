package config

import (
	"easyinvesting/pkg/types"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var secretKey = ""

func SecretKey() string {
	if secretKey == "" {
		secretKey = os.Getenv("EASYINVESTING_SECRET_KEY")
		isProd := os.Getenv("EASYINVESTING_PROD") == "true"
		if !isProd && secretKey == "" {
			secretKey = "easyinvest"
		}
		if secretKey == "" {
			panic("EASYINVESTING_SECRET_KEY IS NOT SET")
		}
	}
	return secretKey
}

func JWTMiddleware() echo.MiddlewareFunc {
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JWTClaims)
		},
		SigningKey: []byte(SecretKey()),
	}
	return echojwt.WithConfig(jwtConfig)
}
