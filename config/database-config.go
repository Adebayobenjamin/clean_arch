package config

import (
	"fmt"
	"os"

	"github.com/Adebayobenjamin/clean_arch/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// create a new connection to our database
func SetupDatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if (err) != nil {
		panic(fmt.Sprintf("Filed to connect to database: %v", err))
	}

	db.AutoMigrate(&entity.Book{}, &entity.User{})
	return db
}

// close connection to database
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Filed to close connection to database: %v", err))
	}
	dbSQL.Close()
}
