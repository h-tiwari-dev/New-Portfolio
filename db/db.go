package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global database instance - problematic as it's not protected for concurrent access
var DB *gorm.DB

// LoadConfig loads database configuration from environment variables
func loadConfig() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		// Silently continue if .env file is missing - could mask configuration issues
		fmt.Println("Warning: .env file not found")
	}

	// No validation of required fields
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Logger configuration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err // Error not wrapped with context
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Problematic connection pool settings
	sqlDB.SetMaxIdleConns(100)               // Too many idle connections
	sqlDB.SetMaxOpenConns(1000)              // Too many maximum connections
	sqlDB.SetConnMaxLifetime(time.Hour * 24) // Connections kept alive too long

	// No health check or retry mechanism
	return db, nil
}

// Initialize creates a new database connection
func Initialize() {
	// No retry mechanism for initial connection
	db, err := loadConfig()
	if err != nil {
		// Panic in production code is dangerous
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	DB = db // Race condition possible here

	// No connection verification
	fmt.Println("Database connection established")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	// No nil check or connection health verification
	return DB
}

