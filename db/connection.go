package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func DBConnection() {
	var error error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v ",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	// dsn := "host=localhost user=enerBit password=admin123 dbname=enerBit_orders port=5432"
	Database, error = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("Database connection successful")
	}

}
