package middleware

import (
	"BE-1/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GrantAdminUmum(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != util.ADMIN ||
			claims["bagian"].(string) != util.UMUM {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
