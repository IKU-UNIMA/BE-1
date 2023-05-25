package handler

import (
	"BE-1/src/api/request"
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/model"
	"BE-1/src/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type alumniQueryParam struct {
	Prodi      int    `query:"prodi"`
	TahunLulus int    `query:"tahun_lulus"`
	Nim        string `query:"nim"`
	Nama       string `query:"nama"`
	Page       int    `query:"page"`
}

func GetAllAlumniHandler(c echo.Context) error {
	queryParams := &alumniQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Alumni{}
	limit := 20
	order := "tahun_lulus DESC"
	conds := ""

	if queryParams.Nim != "" {
		conds = "nim = " + queryParams.Nim
	} else {
		if queryParams.Prodi != 0 {
			conds = fmt.Sprintf("id_prodi = %d", queryParams.Prodi)
		}

		if queryParams.TahunLulus != 0 {
			order = ""
			if conds != "" {
				conds += fmt.Sprintf(" AND tahun_lulus = %d", queryParams.TahunLulus)
			} else {
				conds = fmt.Sprintf("tahun_lulus = %d", queryParams.TahunLulus)
			}
		}

		if queryParams.Nama != "" {
			if conds != "" {
				conds += " AND UPPER(alumni.nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			} else {
				conds = "UPPER(alumni.nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			}
		}
	}

	if err := db.WithContext(ctx).Preload("Prodi").Where(conds).Order(order).
		Offset(util.CountOffset(queryParams.Page, limit)).Limit(limit).
		Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalResult int64
	if err := db.WithContext(ctx).Table("alumni").Where(conds).Count(&totalResult).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Limit:       limit,
		Page:        queryParams.Page,
		TotalPage:   util.CountTotalPage(int(totalResult), limit),
		TotalResult: int(totalResult),
		Data:        data,
	})
}

func GetAlumniByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Alumni{}

	if err := db.WithContext(ctx).Preload("Prodi").Where("id", id).First(data).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "data alumni tidak ditemukan"})
		}
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func ImportAlumniHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []model.Alumni{}

	fileName := file.Filename + time.Now().String()

	defer func() {
		os.Remove(fileName)
	}()

	if err := util.WriteFile(file, fileName); err != nil {
		return err
	}

	excel, err := excelize.OpenFile(fileName)
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer excel.Close()

	rows, err := excel.GetRows(excel.GetSheetName(0))
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if len(rows[0]) != 5 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "jumlah kolom tidak sesuai format"})
	}

	for i := 1; i < len(rows); i++ {
		kodeProdi, err := strconv.Atoi(rows[i][0])
		if err != nil {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("kode prodi pada baris ke-%d tidak valid", i)})
		}

		idProdi := 0
		if err := db.WithContext(ctx).Table("prodi").First(&idProdi, "kode_prodi", kodeProdi).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				message := fmt.Sprintf("prodi dengan kode %d pada baris ke-%d tidak ditemukan", kodeProdi, i)
				return util.FailedResponse(http.StatusNotFound, map[string]string{"message": message})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		tahunLulus, err := strconv.ParseUint(rows[i][4], 10, 32)
		if err != nil {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("tahun lulus pada baris ke-%d tidak valid", i)})
		}

		data = append(data, model.Alumni{
			IdProdi:    idProdi,
			KodePt:     util.KODE_PT,
			Nim:        rows[i][1],
			Nama:       rows[i][2],
			Hp:         rows[i][3],
			TahunLulus: uint(tahunLulus),
		})
	}

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		if strings.Contains(err.Error(), "nim") {
			return util.FailedResponse(
				http.StatusBadRequest,
				map[string]string{"message": "terdapat duplikasi untuk NIM " + strings.Split(err.Error(), "'")[1]},
			)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func InsertAlumniHandler(c echo.Context) error {
	req := &request.InsertAlumni{}
	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := db.WithContext(ctx).Create(req.MapRequest()).Error; err != nil {
		return checkAlumniDBError(err.Error())
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func EditAlumniHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	req := &request.EditAlumni{}
	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := db.WithContext(ctx).Where("id", id).Updates(req.MapRequest()).Error; err != nil {
		return checkAlumniDBError(err.Error())
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteAlumniHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Alumni), id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func checkAlumniDBError(err string) error {
	if strings.Contains(err, "prodi") {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "prodi tidak ditemukan"})
	}

	if strings.Contains(err, "nim") {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nim sudah digunakan"})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
