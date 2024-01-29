package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB SERVER (AWS)
func dbvar() string {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		return ""
	}

	// CONFIG VARS
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")

	dsn := ("host=" + DB_HOST +
		" user=" + DB_USER +
		" password=" + DB_PASSWORD +
		" dbname=" + DB_NAME +
		" port=" + DB_PORT +
		" sslmode=disable TimeZone=Asia/Shanghai") //

	return dsn
}

var DB = func() (db *gorm.DB) {
	if db, err := gorm.Open(postgres.Open(dbvar()), &gorm.Config{}); err != nil {
		fmt.Println("Connection to database failed", err)
		panic(err)
	} else {
		fmt.Println("Connected to database")
		return db
	}
}()
