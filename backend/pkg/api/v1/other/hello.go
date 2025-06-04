package other

import (
	"easyinvesting/pkg/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, types.JsonMap{"message": "Hello from Go Echo!"})
}
