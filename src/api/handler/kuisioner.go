package handler

import (
	"BE-1/src/api/request"
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/model"
	"BE-1/src/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type kuisionerQueryParam struct {
	Prodi      int    `query:"prodi"`
	Nim        string `query:"nim"`
	Nama       string `query:"nama"`
	TahunLulus int    `query:"tahun_lulus"`
	Limit      int    `query:"limit"`
	Page       int    `query:"page"`
}

func CheckKuisionerByNIMHandler(c echo.Context) error {
	nim := c.Param("nim")

	db := database.InitMySQL()
	ctx := c.Request().Context()
	alumni := &response.Alumni{}

	if err := db.WithContext(ctx).Preload("Prodi").First(&alumni, "nim", nim).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "alumni tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	id := 0
	if err := db.WithContext(ctx).Table("kuisioner").Select("id").Where("id_alumni", alumni.ID).Scan(&id).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if id != 0 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "alumni sudah mengisi kuisioner"})
	}

	return util.SuccessResponse(c, http.StatusOK, alumni)
}

func GetAllKuisionerHandler(c echo.Context) error {
	queryParams := &kuisionerQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	conds := ""
	if queryParams.Nim != "" {
		conds = "alumni.nim = " + queryParams.Nim
	} else {
		if queryParams.Prodi != 0 {
			conds = fmt.Sprintf("alumni.id_prodi = %d", queryParams.Prodi)
		}

		if queryParams.TahunLulus != 0 {
			if conds != "" {
				conds += fmt.Sprintf(" AND alumni.tahun_lulus = %d", queryParams.TahunLulus)
			} else {
				conds = fmt.Sprintf("alumni.tahun_lulus = %d", queryParams.TahunLulus)
			}
		}

		if queryParams.Nama != "" {
			if conds != "" {
				conds += " AND UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			} else {
				conds = "UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			}
		}
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Kuisioner{}

	if err := db.WithContext(ctx).Preload("Alumni.Prodi").
		Joins("JOIN alumni ON alumni.id = kuisioner.id_alumni").
		Where(conds).
		Offset(util.CountOffset(queryParams.Page, queryParams.Limit)).Limit(queryParams.Limit).
		Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalResult int64
	if err := db.WithContext(ctx).Table("kuisioner").
		Joins("JOIN alumni ON alumni.id = kuisioner.id_alumni").
		Where(conds).Count(&totalResult).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetKuisionerByIDHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DetailKuisioner{}

	if err := db.WithContext(ctx).Table("kuisioner").Preload("Alumni.Prodi").First(&data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func ImportKuisionerHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()
	data := []model.Kuisioner{}

	fileName := util.GetNewFileName(file.Filename)

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

	for i := 1; i < len(rows); i++ {
		kdProdi := rows[i][1]
		idProdi := 0
		if err := db.WithContext(ctx).Model(new(model.Prodi)).Select("id").First(&idProdi, "kode_prodi", kdProdi).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				message := fmt.Sprintf("prodi dengan kode %s pada baris ke-%d tidak ditemukan", kdProdi, i)
				return util.FailedResponse(http.StatusNotFound, map[string]string{"message": message})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		tahunLulus, err := strconv.ParseUint(rows[i][6], 10, 32)
		if err != nil {
			message := fmt.Sprintf("tahun lulus pada baris ke-%d tidak valid", i)
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
		}

		nik, err := strconv.Atoi(rows[i][7])
		if err != nil {
			message := fmt.Sprintf("NIK pada baris ke-%d tidak valid", i)
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
		}

		alumni := model.Alumni{
			IdProdi:    idProdi,
			KodePt:     util.KODE_PT,
			Nim:        rows[i][2],
			Nama:       rows[i][3],
			Hp:         rows[i][4],
			Email:      &rows[i][5],
			TahunLulus: uint(tahunLulus),
			Nik:        &nik,
			Npwp:       &rows[i][8],
		}

		// find alumni
		idAlumni := 0
		if err := db.WithContext(ctx).Table("alumni").Select("id").
			Where("nim", alumni.Nim).Scan(&idAlumni).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		if idAlumni != 0 {
			message := fmt.Sprintf("alumni dengan nim %s pada baris ke-%d sudah mengisi kuisioner", alumni.Nim, i)
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
		}

		if err := tx.WithContext(ctx).Where("nim", alumni.Nim).Create(&alumni).Error; err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "nik") {
				message := fmt.Sprintf("NIK %d pada baris ke-%d sudah digunakan", *alumni.Nik, i)
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
			}

			if strings.Contains(err.Error(), "nik") {
				message := fmt.Sprintf("NPWP %s pada baris ke-%d sudah digunakan", *alumni.Npwp, i)
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		idAlumni = alumni.ID

		f8, _ := strconv.Atoi(rows[i][9])
		f504, _ := strconv.Atoi(rows[i][10])
		f502, _ := strconv.Atoi(rows[i][11])
		f505, _ := strconv.Atoi(rows[i][12])
		f5a1, _ := strconv.Atoi(rows[i][14])
		f5a2, _ := strconv.Atoi(rows[i][15])
		f1101, _ := strconv.Atoi(rows[i][16])
		f5c, _ := strconv.Atoi(rows[i][19])
		f5d, _ := strconv.Atoi(rows[i][20])
		f18a, _ := strconv.Atoi(rows[i][21])
		f1201, _ := strconv.Atoi(rows[i][25])
		f14, _ := strconv.Atoi(rows[i][27])
		f15, _ := strconv.Atoi(rows[i][28])
		f1761, _ := strconv.Atoi(rows[i][29])
		f1762, _ := strconv.Atoi(rows[i][30])
		f1763, _ := strconv.Atoi(rows[i][31])
		f1764, _ := strconv.Atoi(rows[i][32])
		f1765, _ := strconv.Atoi(rows[i][33])
		f1766, _ := strconv.Atoi(rows[i][34])
		f1767, _ := strconv.Atoi(rows[i][35])
		f1768, _ := strconv.Atoi(rows[i][36])
		f1769, _ := strconv.Atoi(rows[i][37])
		f1770, _ := strconv.Atoi(rows[i][38])
		f1771, _ := strconv.Atoi(rows[i][39])
		f1772, _ := strconv.Atoi(rows[i][40])
		f1773, _ := strconv.Atoi(rows[i][41])
		f1774, _ := strconv.Atoi(rows[i][42])
		f21, _ := strconv.Atoi(rows[i][43])
		f22, _ := strconv.Atoi(rows[i][44])
		f23, _ := strconv.Atoi(rows[i][45])
		f24, _ := strconv.Atoi(rows[i][46])
		f25, _ := strconv.Atoi(rows[i][47])
		f26, _ := strconv.Atoi(rows[i][48])
		f27, _ := strconv.Atoi(rows[i][49])
		f301, _ := strconv.Atoi(rows[i][50])
		f302, _ := strconv.Atoi(rows[i][51])
		f303, _ := strconv.Atoi(rows[i][52])
		f401, _ := strconv.Atoi(rows[i][53])
		f402, _ := strconv.Atoi(rows[i][54])
		f403, _ := strconv.Atoi(rows[i][55])
		f404, _ := strconv.Atoi(rows[i][56])
		f405, _ := strconv.Atoi(rows[i][57])
		f406, _ := strconv.Atoi(rows[i][58])
		f407, _ := strconv.Atoi(rows[i][59])
		f408, _ := strconv.Atoi(rows[i][60])
		f409, _ := strconv.Atoi(rows[i][61])
		f410, _ := strconv.Atoi(rows[i][62])
		f411, _ := strconv.Atoi(rows[i][63])
		f412, _ := strconv.Atoi(rows[i][64])
		f413, _ := strconv.Atoi(rows[i][65])
		f414, _ := strconv.Atoi(rows[i][66])
		f415, _ := strconv.Atoi(rows[i][67])
		f6, _ := strconv.Atoi(rows[i][69])
		f7, _ := strconv.Atoi(rows[i][70])
		f7a, _ := strconv.Atoi(rows[i][71])
		f1001, _ := strconv.Atoi(rows[i][72])
		f1601, _ := strconv.Atoi(rows[i][74])
		f1602, _ := strconv.Atoi(rows[i][75])
		f1603, _ := strconv.Atoi(rows[i][76])
		f1604, _ := strconv.Atoi(rows[i][77])
		f1605, _ := strconv.Atoi(rows[i][78])
		f1606, _ := strconv.Atoi(rows[i][79])
		f1607, _ := strconv.Atoi(rows[i][80])
		f1608, _ := strconv.Atoi(rows[i][81])
		f1609, _ := strconv.Atoi(rows[i][82])
		f1610, _ := strconv.Atoi(rows[i][83])
		f1611, _ := strconv.Atoi(rows[i][84])
		f1612, _ := strconv.Atoi(rows[i][85])
		f1613, _ := strconv.Atoi(rows[i][86])

		kuisioner := &request.EditKuisioner{
			F8:    int8(f8),
			F504:  int8(f504),
			F502:  int8(f502),
			F505:  int32(f505),
			F5a1:  int32(f5a1),
			F5a2:  int32(f5a2),
			F1101: int8(f1101),
			F1102: rows[i][17],
			F5b:   rows[i][18],
			F5c:   int8(f5c),
			F5d:   int8(f5d),
			F18a:  int8(f18a),
			F18b:  rows[i][22],
			F18c:  rows[i][23],
			F18d:  rows[i][24],
			F1201: int8(f1201),
			F1202: rows[i][26],
			F14:   int8(f14),
			F15:   int8(f15),
			F1761: int8(f1761),
			F1762: int8(f1762),
			F1763: int8(f1763),
			F1764: int8(f1764),
			F1765: int8(f1765),
			F1766: int8(f1766),
			F1767: int8(f1767),
			F1768: int8(f1768),
			F1769: int8(f1769),
			F1770: int8(f1770),
			F1771: int8(f1771),
			F1772: int8(f1772),
			F1773: int8(f1773),
			F1774: int8(f1774),
			F21:   int8(f21),
			F22:   int8(f22),
			F23:   int8(f23),
			F24:   int8(f24),
			F25:   int8(f25),
			F26:   int8(f26),
			F27:   int8(f27),
			F301:  int8(f301),
			F302:  int8(f302),
			F303:  int8(f303),
			F401:  int8(f401),
			F402:  int8(f402),
			F403:  int8(f403),
			F404:  int8(f404),
			F405:  int8(f405),
			F406:  int8(f406),
			F407:  int8(f407),
			F408:  int8(f408),
			F409:  int8(f409),
			F410:  int8(f410),
			F411:  int8(f411),
			F412:  int8(f412),
			F413:  int8(f413),
			F414:  int8(f414),
			F415:  int8(f415),
			F416:  rows[i][68],
			F6:    int16(f6),
			F7:    int16(f7),
			F7a:   int16(f7a),
			F1001: int8(f1001),
			F1002: rows[i][73],
			F1601: int8(f1601),
			F1602: int8(f1602),
			F1603: int8(f1603),
			F1604: int8(f1604),
			F1605: int8(f1605),
			F1606: int8(f1606),
			F1607: int8(f1607),
			F1608: int8(f1608),
			F1609: int8(f1609),
			F1610: int8(f1610),
			F1611: int8(f1611),
			F1612: int8(f1612),
			F1613: int8(f1613),
			F1614: rows[i][87],
		}

		data = append(data, *kuisioner.MapRequest())
		data[i-1].IdAlumni = idAlumni

		if err := tx.WithContext(ctx).Create(&data[i-1]).Error; err != nil {
			tx.Rollback()
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ExportKuisionerHandler(c echo.Context) error {
	queryParams := &kuisionerQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	conds := ""
	if queryParams.Prodi != 0 {
		conds = fmt.Sprintf("alumni.id_prodi = %d", queryParams.Prodi)
	}

	if queryParams.TahunLulus != 0 {
		if conds != "" {
			conds += fmt.Sprintf(" AND alumni.tahun_lulus = %d", queryParams.TahunLulus)
		} else {
			conds = fmt.Sprintf("alumni.tahun_lulus = %d", queryParams.TahunLulus)
		}
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DetailKuisioner{}

	if err := db.WithContext(ctx).Table("kuisioner").Preload("Alumni.Prodi").
		Joins("JOIN alumni ON alumni.id = kuisioner.id_alumni").
		Where(conds).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	fileName := "data_responden.xlsx"
	newFile := util.GetNewFileName(fileName)

	defer func() {
		os.Remove(newFile)
	}()

	src, err := os.Open(fileName)
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer src.Close()

	dst, err := os.Create(newFile)
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer f.Close()

	for i := 0; i < len(data); i++ {
		sheet := "upload"
		createCell := func(cell string) string {
			return fmt.Sprintf("%s%d", cell, i+2)
		}

		f.SetCellValue(sheet, createCell("A"), util.KODE_PT)
		f.SetCellValue(sheet, createCell("B"), data[i].Alumni.Prodi.KodeProdi)
		f.SetCellValue(sheet, createCell("C"), data[i].Alumni.Nim)
		f.SetCellValue(sheet, createCell("D"), data[i].Alumni.Nama)
		f.SetCellValue(sheet, createCell("E"), data[i].Alumni.Hp)
		f.SetCellValue(sheet, createCell("F"), data[i].Alumni.Email)
		f.SetCellValue(sheet, createCell("G"), data[i].Alumni.TahunLulus)
		f.SetCellValue(sheet, createCell("H"), data[i].Alumni.Nik)
		f.SetCellValue(sheet, createCell("I"), data[i].Alumni.Npwp)
		f.SetCellValue(sheet, createCell("J"), data[i].F8)
		f.SetCellValue(sheet, createCell("K"), data[i].F504)
		f.SetCellValue(sheet, createCell("L"), data[i].F502)
		f.SetCellValue(sheet, createCell("M"), data[i].F505)
		f.SetCellValue(sheet, createCell("N"), "")
		f.SetCellValue(sheet, createCell("O"), data[i].F5a1)
		f.SetCellValue(sheet, createCell("P"), data[i].F5a2)
		f.SetCellValue(sheet, createCell("Q"), data[i].F1101)
		f.SetCellValue(sheet, createCell("R"), data[i].F1102)
		f.SetCellValue(sheet, createCell("S"), data[i].F5b)
		f.SetCellValue(sheet, createCell("T"), data[i].F5c)
		f.SetCellValue(sheet, createCell("U"), data[i].F5d)
		f.SetCellValue(sheet, createCell("V"), data[i].F18a)
		f.SetCellValue(sheet, createCell("W"), data[i].F18b)
		f.SetCellValue(sheet, createCell("X"), data[i].F18c)
		f.SetCellValue(sheet, createCell("Y"), data[i].F18d)
		f.SetCellValue(sheet, createCell("Z"), data[i].F1201)
		f.SetCellValue(sheet, createCell("AA"), data[i].F1202)
		f.SetCellValue(sheet, createCell("AB"), data[i].F14)
		f.SetCellValue(sheet, createCell("AC"), data[i].F15)
		f.SetCellValue(sheet, createCell("AD"), data[i].F1761)
		f.SetCellValue(sheet, createCell("AE"), data[i].F1762)
		f.SetCellValue(sheet, createCell("AF"), data[i].F1763)
		f.SetCellValue(sheet, createCell("AG"), data[i].F1764)
		f.SetCellValue(sheet, createCell("AH"), data[i].F1765)
		f.SetCellValue(sheet, createCell("AI"), data[i].F1766)
		f.SetCellValue(sheet, createCell("AJ"), data[i].F1767)
		f.SetCellValue(sheet, createCell("AK"), data[i].F1768)
		f.SetCellValue(sheet, createCell("AL"), data[i].F1769)
		f.SetCellValue(sheet, createCell("AM"), data[i].F1770)
		f.SetCellValue(sheet, createCell("AN"), data[i].F1771)
		f.SetCellValue(sheet, createCell("AO"), data[i].F1772)
		f.SetCellValue(sheet, createCell("AP"), data[i].F1773)
		f.SetCellValue(sheet, createCell("AQ"), data[i].F1774)
		f.SetCellValue(sheet, createCell("AR"), data[i].F21)
		f.SetCellValue(sheet, createCell("AS"), data[i].F22)
		f.SetCellValue(sheet, createCell("AT"), data[i].F23)
		f.SetCellValue(sheet, createCell("AU"), data[i].F24)
		f.SetCellValue(sheet, createCell("AV"), data[i].F25)
		f.SetCellValue(sheet, createCell("AW"), data[i].F26)
		f.SetCellValue(sheet, createCell("AX"), data[i].F27)
		f.SetCellValue(sheet, createCell("AY"), data[i].F301)
		f.SetCellValue(sheet, createCell("AZ"), data[i].F302)
		f.SetCellValue(sheet, createCell("BA"), data[i].F303)
		f.SetCellValue(sheet, createCell("BB"), data[i].F401)
		f.SetCellValue(sheet, createCell("BC"), data[i].F402)
		f.SetCellValue(sheet, createCell("BD"), data[i].F403)
		f.SetCellValue(sheet, createCell("BE"), data[i].F404)
		f.SetCellValue(sheet, createCell("BF"), data[i].F405)
		f.SetCellValue(sheet, createCell("BG"), data[i].F406)
		f.SetCellValue(sheet, createCell("BH"), data[i].F407)
		f.SetCellValue(sheet, createCell("BI"), data[i].F408)
		f.SetCellValue(sheet, createCell("BJ"), data[i].F409)
		f.SetCellValue(sheet, createCell("BK"), data[i].F410)
		f.SetCellValue(sheet, createCell("BL"), data[i].F411)
		f.SetCellValue(sheet, createCell("BM"), data[i].F412)
		f.SetCellValue(sheet, createCell("BN"), data[i].F413)
		f.SetCellValue(sheet, createCell("BO"), data[i].F414)
		f.SetCellValue(sheet, createCell("BP"), data[i].F415)
		f.SetCellValue(sheet, createCell("BQ"), data[i].F416)
		f.SetCellValue(sheet, createCell("BR"), data[i].F6)
		f.SetCellValue(sheet, createCell("BS"), data[i].F7)
		f.SetCellValue(sheet, createCell("BT"), data[i].F7a)
		f.SetCellValue(sheet, createCell("BU"), data[i].F1001)
		f.SetCellValue(sheet, createCell("BV"), data[i].F1002)
		f.SetCellValue(sheet, createCell("BW"), data[i].F1601)
		f.SetCellValue(sheet, createCell("BX"), data[i].F1602)
		f.SetCellValue(sheet, createCell("BY"), data[i].F1603)
		f.SetCellValue(sheet, createCell("BZ"), data[i].F1604)
		f.SetCellValue(sheet, createCell("CA"), data[i].F1605)
		f.SetCellValue(sheet, createCell("CB"), data[i].F1606)
		f.SetCellValue(sheet, createCell("CC"), data[i].F1607)
		f.SetCellValue(sheet, createCell("CD"), data[i].F1608)
		f.SetCellValue(sheet, createCell("CE"), data[i].F1609)
		f.SetCellValue(sheet, createCell("CF"), data[i].F1610)
		f.SetCellValue(sheet, createCell("CG"), data[i].F1611)
		f.SetCellValue(sheet, createCell("CH"), data[i].F1612)
		f.SetCellValue(sheet, createCell("CI"), data[i].F1613)
		f.SetCellValue(sheet, createCell("CJ"), data[i].F1614)
	}

	if err := f.Save(); err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return c.File(newFile)
}

func InsertKuisionerHandler(c echo.Context) error {
	req := &request.Kuisioner{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	idAlumni := 0
	if err := db.WithContext(ctx).Model(new(model.Alumni)).Select("id").First(&idAlumni).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "alumni tidak ditemukan"})
		}

		util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Create(req.MapRequest(idAlumni)).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "alumni sudah mengisi kuisioner"})
		}

		if strings.Contains(err.Error(), "alumni_nim") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nim tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Where("nim", req.Nim).Updates(req.MapAlumniData()).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "npwp") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NPWP sudah digunakan"})
		}

		if strings.Contains(err.Error(), "nik") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIK sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditKuisionerHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	req := &request.EditKuisioner{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).Where("id", id).Updates(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKuisionerHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Kuisioner), "id", id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ApproveKuisionerHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).Where("id", id).Update("status", true).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
