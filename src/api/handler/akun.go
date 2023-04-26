package handler

import (
	"BE-1/src/api/request"
	"BE-1/src/api/response"
	"BE-1/src/config/database"
	"BE-1/src/model"
	"BE-1/src/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	req := &request.Login{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}

	if err := db.WithContext(ctx).First(data, "email", req.Email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(req.Password, data.Password) {
		return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
	}

	var bagian string
	if data.Role == string(util.ADMIN) {
		if err := db.WithContext(ctx).Table("admin").Select("bagian").Where("id", data.ID).Scan(&bagian).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	var nama string
	if err := db.WithContext(ctx).Table(data.Role).Select("nama").Where("id", data.ID).Scan(&nama).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	token := util.GenerateJWT(data.ID, nama, data.Role, bagian)

	return util.SuccessResponse(c, http.StatusOK, response.Login{Token: token})
}

func CheckNIMHandler(c echo.Context) error {
	req := &request.CheckNIM{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	alumni := struct {
		ID  int
		Nim string
	}{}

	if err := db.WithContext(ctx).Table("alumni").Select("id", "nim").Where("nim", req.Nim).Scan(&alumni).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if alumni.Nim != req.Nim {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "data tidak ditemukan"})
	}

	conds := fmt.Sprintf("id = %d AND password IS NULL", alumni.ID)
	if err := db.WithContext(ctx).First(new(model.Akun), conds).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "akun sudah terdaftar"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, req)
}

func RegisterAlumniHandler(c echo.Context) error {
	req := &request.RegisterAlumni{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	alumni := struct {
		ID  int
		Nim string
	}{}

	if err := db.WithContext(ctx).Table("alumni").Select("id", "nim").Where("nim", req.Nim).Scan(&alumni).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if alumni.Nim != req.Nim {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "data tidak ditemukan"})
	}

	conds := fmt.Sprintf("id = %d AND password IS NULL", alumni.ID)
	if err := db.WithContext(ctx).First(new(model.Akun), conds).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "akun sudah terdaftar"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	akun := req.MapRequest()

	query := fmt.Sprintf(`
		UPDATE akun
		SET 
			akun.email = '%s',
			akun.password = '%s'
		WHERE akun.id = %d;
	`, akun.Email, akun.Password, alumni.ID)

	if err := db.WithContext(ctx).Exec(query).Error; err != nil {
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ChangePasswordHandler(c echo.Context) error {
	req := &request.ChangePassword{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(req.PasswordBaru)).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
