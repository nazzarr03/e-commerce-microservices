package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nazzarr03/product-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ConnectProductDB()
}

func ConnectProductDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(&models.Product{}); err != nil {
		panic("failed to migrate database")
	}

	log.Println("Successfully connected to database")
}
