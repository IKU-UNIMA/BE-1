package handler

import (
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllFakultasHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	result := []response.Fakultas{}

	if err := db.WithContext(ctx).Order("id").Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}
