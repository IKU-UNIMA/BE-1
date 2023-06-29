package database

import (
	"BE-1/src/config/env"
	"BE-1/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(env.GetMySQLEnv()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
}

func MigrateMySQL() {
	DB.AutoMigrate(
		&model.Provinsi{},
		&model.KabKota{},
		&model.Fakultas{},
		&model.Prodi{},
		&model.Kuisioner{},
		&model.Alumni{},
		&model.Akun{},
		&model.Admin{},
		&model.Rektor{},
		&model.Target{},
	)
}
