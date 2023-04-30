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

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type alumniQueryParam struct {
	Prodi      int    `query:"prodi"`
	TahunLulus int    `query:"tahun_lulus"`
	Nim        string `query:"nim"`
	Nama       string `query:"nama"`
}

const getAlumniQuery = `
	SELECT alumni.id, prodi.nama AS prodi, nim, alumni.nama, akun.email, hp, tahun_lulus, npwp, nik 
	FROM alumni 
	JOIN akun ON alumni.id = akun.id 
	JOIN prodi ON prodi.id = alumni.id_prodi
	`

func GetAllAlumniHandler(c echo.Context) error {
	queryParams := &alumniQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Alumni{}
	conds := ""

	if queryParams.Nim != "" {
		conds = "nim = " + queryParams.Nim
	} else {
		if queryParams.Prodi != 0 {
			conds = fmt.Sprintf("id_prodi = %d", queryParams.Prodi)
		}

		if queryParams.TahunLulus != 0 {
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

	if conds != "" {
		conds = " WHERE " + conds
	}

	query := getAlumniQuery + conds
	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetAlumniByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Alumni{}

	query := getAlumniQuery + fmt.Sprintf(" WHERE alumni.id = %d", id)
	if err := db.WithContext(ctx).Raw(query).First(data).Error; err != nil {
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

	defer func() {
		os.Remove(file.Filename)
	}()

	if err := util.WriteFile(file); err != nil {
		return err
	}

	excel, err := excelize.OpenFile(file.Filename)
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
		idProdi, err := strconv.Atoi(rows[i][0])
		if err != nil {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("kode prodi pada baris ke-%d tidak valid", i)})
		}

		tahunLulus, err := strconv.ParseUint(rows[i][4], 10, 32)
		if err != nil {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("tahun lulus pada baris ke-%d tidak valid", i)})
		}

		data = append(data, model.Alumni{
			IdProdi:    idProdi,
			KodePt:     "001035",
			Nim:        rows[i][1],
			Nama:       rows[i][2],
			Hp:         rows[i][3],
			TahunLulus: uint(tahunLulus),
			Akun: model.Akun{
				Role: util.ALUMNI,
			},
		})
	}

	if err := db.WithContext(ctx).Omit("Akun.Email", "Akun.Password").Create(&data).Error; err != nil {
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
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Omit("Akun.Email", "Akun.Password").Create(req.MapRequest()).Error; err != nil {
		tx.Rollback()
		return checkAlumniDBError(err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
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
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Where("id", id).Updates(req.MapRequest()).Error; err != nil {
		tx.Rollback()
		return checkAlumniDBError(err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
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
	if strings.Contains(err, "nim") {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nim sudah digunakan"})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
