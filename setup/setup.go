package setup

import (
	"backend_golang/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Database configuration variables are not set properly.")
	}
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automigrate your models
	err = database.AutoMigrate(
		&models.User{},
		&models.Siswa{},
		&models.Menu{},
		&models.Jenis{},
		&models.Stan{},
		&models.Status{},
		&models.Detail{},
		&models.Transaksi{},
		&models.Role{},
		&models.Diskon{},
		&models.MenuDiskon{},

	)
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	DB = database
}
