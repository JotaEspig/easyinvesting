package other

import (
	"easyinvesting/pkg/api/v1/utils"
	"easyinvesting/pkg/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func helloAuth(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.JsonMap{"error": "Unauthorized"})
	}

	user := types.JsonMap{
		"id": claims.UserID,
	}
	return c.JSON(http.StatusOK, types.JsonMap{"message": "Hello from Go Echo with Auth!", "user": user})
}
