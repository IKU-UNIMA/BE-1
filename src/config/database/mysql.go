package database

import (
	"BE-1/src/config/env"
	"BE-1/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitMySQL() *gorm.DB {
	db, err := gorm.Open(mysql.Open(env.GetMySQLEnv()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	return db
}

func MigrateMySQL() {
	InitMySQL().AutoMigrate(
		&model.Provinsi{},
		&model.KabKota{},
		&model.Prodi{},
		&model.Admin{},
		&model.Alumni{},
	)
}
