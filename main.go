package main

import (
	"BE-1/src/api/route"
	"BE-1/src/config/database"
	"BE-1/src/config/env"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init mysql
	database.InitMySQL()

	// migrate gorm
	database.MigrateMySQL()

	app := route.InitServer()
	app.Logger.Fatal(app.Start(":" + env.GetServerEnv()))
}
