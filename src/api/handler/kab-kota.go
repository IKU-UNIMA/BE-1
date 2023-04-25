package handler

import (
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllKabKotaByProvinsi(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.KabKota{}
	idProvinsi := c.Param("id")

	if err := db.WithContext(ctx).Find(&data, "id_provinsi", idProvinsi).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}
