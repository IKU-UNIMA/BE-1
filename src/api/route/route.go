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

	provinsi := v1.Group("/provinsi", customMiddleware.Authentication)
	provinsi.GET("", handler.GetAllProvinsiHandler)

	kabKota := v1.Group("/kab-kota", customMiddleware.Authentication)
	kabKota.GET("/provinsi/:id", handler.GetAllKabKotaByProvinsi)

	fakultas := v1.Group("/fakultas", customMiddleware.Authentication)
	fakultas.GET("", handler.GetAllFakultasHandler)

	prodi := v1.Group("/prodi", customMiddleware.Authentication)
	prodi.GET("", handler.GetAllProdiHandler)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.POST("/register/alumni", handler.RegisterAlumniHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)

	admin := v1.Group("/admin", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	admin.GET("", handler.GetAllAdminHandler)
	admin.GET("/:id", handler.GetAdminByIdHandler)
	admin.POST("", handler.InsertAdminHandler)
	admin.PUT("/:id", handler.EditAdminHandler)
	admin.DELETE("/:id", handler.DeleteAdminHandler)

	alumni := v1.Group("/alumni", customMiddleware.Authentication, customMiddleware.GrantAdminIKU1)
	alumni.GET("", handler.GetAllAlumniHandler)
	alumni.GET("/belum-mengisi", handler.GetAllAlumniBelumMengisiHandler)
	alumni.GET("/:id", handler.GetAlumniByIdHandler)
	alumni.POST("", handler.InsertAlumniHandler)
	alumni.POST("/import", handler.ImportAlumniHandler)
	alumni.PUT("/:id", handler.EditAlumniHandler)
	alumni.DELETE("/:id", handler.DeleteAlumniHandler)

	rektor := v1.Group("/rektor", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	rektor.GET("", handler.GetAllRektorHandler)
	rektor.GET("/:id", handler.GetRektorByIdHandler)
	rektor.POST("", handler.InsertRektorHandler)
	rektor.PUT("/:id", handler.EditRektorHandler)
	rektor.DELETE("/:id", handler.DeleteRektorHandler)

	kuisioner := v1.Group("/kuisioner")
	kuisioner.GET("/check/:nim", handler.CheckKuisionerByNIMHandler)
	kuisioner.POST("", handler.InsertKuisionerHandler)

	kuisionerAuth := kuisioner.Group("", customMiddleware.Authentication, customMiddleware.GrantAdminIKU1)
	kuisionerAuth.POST("/import", handler.ImportKuisionerHandler)
	kuisionerAuth.GET("/export", handler.ExportKuisionerHandler)
	kuisionerAuth.GET("", handler.GetAllKuisionerHandler)
	kuisionerAuth.GET("/:id", handler.GetKuisionerByIDHandler)
	kuisionerAuth.PUT("/:id", handler.EditKuisionerHandler)
	kuisionerAuth.DELETE("/:id", handler.DeleteKuisionerHandler)
	kuisionerAuth.PATCH("/:id/approve", handler.ApproveKuisionerHandler)

	dashboard := v1.Group("/dashboard", customMiddleware.Authentication)
	dashboard.GET("/tahun/:tahun", handler.GetDashboardHandler, customMiddleware.GrantAdminIKU1AndRektor)
	dashboard.GET("/fakultas/:fakultas/:tahun", handler.GetDashboardByFakultasHandler, customMiddleware.GrantAdminIKU1AndRektor)
	dashboard.PATCH("/target", handler.InsertTargetHandler, customMiddleware.GrantAdminIKU1)

	return app
}
