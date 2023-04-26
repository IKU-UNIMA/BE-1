package route

import (
	"BE-1/src/api/handler"
	"BE-1/src/util/validation"

	customMiddleware "BE-1/src/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {
	app := echo.New()
	app.Use(middleware.CORS())

	app.Validator = &validation.CustomValidator{Validator: validator.New()}

	app.GET("", func(c echo.Context) error {
		return c.JSON(200, "Welcome to IKU 1 API")
	})

	v1 := app.Group("/api/v1")

	provinsi := v1.Group("/provinsi")
	provinsi.GET("", handler.GetAllProvinsiHandler)

	kabKota := v1.Group("/kab-kota")
	kabKota.GET("/provinsi/:id", handler.GetAllKabKotaByProvinsi)

	prodi := v1.Group("/prodi")
	prodi.GET("", handler.GetAllProdiHandler)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.POST("/register/alumni/check", handler.CheckNIMHandler)
	akun.POST("/register/alumni", handler.RegisterAlumniHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)

	admin := v1.Group("/admin", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	admin.GET("", handler.GetAllAdminHandler)
	admin.GET("/:id", handler.GetAdminByIdHandler)
	admin.POST("", handler.InsertAdminHandler)
	admin.PUT("/:id", handler.EditAdminHandler)
	admin.DELETE("/:id", handler.DeleteAdminHandler)

	return app
}
