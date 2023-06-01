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

func GrantAdminIKU1(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != util.ADMIN ||
			claims["bagian"].(string) != util.IKU1 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU1AndRektor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.REKTOR) && role != string(util.ADMIN) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != util.IKU1 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
