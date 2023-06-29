package handler

import (
	"BE-1/src/api/request"
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardPathParam struct {
	Fakultas int `param:"fakultas"`
	Tahun    int `param:"tahun"`
}

func GetDashboardHandler(c echo.Context) error {
	params := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, params); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &response.Dashboard{}

	var target float64
	targetQuery := fmt.Sprintf(`
	SELECT target FROM target
	WHERE bagian = 'IKU 1' AND tahun = %d
	`, params.Tahun)
	if err := db.WithContext(ctx).Raw(targetQuery).Find(&target).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Target = fmt.Sprintf("%.1f", util.RoundFloat(target))

	conds := fmt.Sprintf("AND alumni.tahun_lulus= %d", params.Tahun)
	query := fmt.Sprintf(`
	SELECT
		fakultas.id, fakultas.nama AS fakultas, COUNT(alumni.id) AS jumlah_alumni,
		COUNT(kuisioner.id) AS jumlah_responden,
		(sum(if(kuisioner.f8 = 1, 1, 0))) AS bekerja,
		(sum(if(kuisioner.f8 = 3, 1, 0))) AS wiraswasta,
		(sum(if(kuisioner.f8 = 4, 1, 0))) AS melanjutkan_pendidikan
	FROM fakultas
	LEFT JOIN prodi ON prodi.id_fakultas = fakultas.id
	LEFT JOIN alumni ON alumni.id_prodi = prodi.id %s
	LEFT JOIN kuisioner ON kuisioner.id_alumni = alumni.id
		AND kuisioner.f8 IN (1, 3, 4)
	GROUP BY fakultas.id ORDER BY fakultas.id
	`, conds)

	if err := db.WithContext(ctx).Raw(query).Find(&data.Detail).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	for i := 0; i < len(data.Detail); i++ {
		data.Total += data.Detail[i].JumlahResponden
		data.TotalAlumni += data.Detail[i].JumlahAlumni

		var persentase float64
		if data.Detail[i].JumlahResponden != 0 {
			persentase = util.RoundFloat((float64(data.Detail[i].JumlahResponden) / float64(data.Detail[i].JumlahAlumni)) * 100)
		}

		data.Detail[i].Persentase = fmt.Sprintf("%.2f", persentase) + "%"
	}

	pencapaian := util.RoundFloat((float64(data.Total) / float64(data.TotalAlumni)) * 100)
	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardByFakultasHandler(c echo.Context) error {
	params := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, params); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &response.DashboardPerProdi{}

	fakultas := ""
	if err := db.WithContext(ctx).Raw("SELECT nama FROM fakultas WHERE id = ?", params.Fakultas).First(&fakultas).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Fakultas = fakultas

	conds := fmt.Sprintf("AND alumni.tahun_lulus= %d", params.Tahun)
	query := fmt.Sprintf(`
	SELECT
		CONCAT(kode_prodi, " - ", prodi.nama, " (", prodi.jenjang, ")") AS prodi,
		COUNT(alumni.id) AS jumlah_alumni,
		COUNT(kuisioner.id) AS jumlah_responden,
		(sum(if(kuisioner.f8 = 1, 1, 0))) AS bekerja,
		(sum(if(kuisioner.f8 = 3, 1, 0))) AS wiraswasta,
		(sum(if(kuisioner.f8 = 4, 1, 0))) AS melanjutkan_pendidikan
	FROM prodi
	LEFT JOIN alumni ON alumni.id_prodi = prodi.id %s
	LEFT JOIN kuisioner ON kuisioner.id_alumni = alumni.id
		AND kuisioner.f8 IN (1, 3, 4)
	WHERE prodi.id_fakultas = %d
	GROUP BY prodi.id ORDER BY prodi.id
	`, conds, params.Fakultas)

	if err := db.WithContext(ctx).Raw(query).Find(&data.Detail).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	for i := 0; i < len(data.Detail); i++ {
		data.Total += data.Detail[i].JumlahResponden
		data.TotalAlumni += data.Detail[i].JumlahAlumni

		var persentase float64
		if data.Detail[i].JumlahResponden != 0 {
			persentase = util.RoundFloat((float64(data.Detail[i].JumlahResponden) / float64(data.Detail[i].JumlahAlumni)) * 100)
		}

		data.Detail[i].Persentase = fmt.Sprintf("%.2f", persentase) + "%"
	}

	pencapaian := util.RoundFloat((float64(data.Total) / float64(data.TotalAlumni)) * 100)
	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertTargetHandler(c echo.Context) error {
	req := &request.Target{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	conds := fmt.Sprintf("bagian='%s' AND tahun=%d", util.IKU1, req.Tahun)

	if err := db.WithContext(ctx).Where(conds).Save(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
