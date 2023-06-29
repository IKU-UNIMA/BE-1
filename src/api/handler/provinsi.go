package handler

import (
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllProvinsiHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	data := []response.Provinsi{}

	if err := db.WithContext(ctx).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}
