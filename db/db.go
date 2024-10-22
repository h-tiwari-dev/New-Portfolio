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

// Config holds database configuration parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		TimeZone: getEnvOrDefault("DB_TIMEZONE", "Asia/Shanghai"),
	}

	if config.Password == "" {
		return nil, fmt.Errorf("database password not set in environment")
	}

	return config, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// BuildDSN constructs the database connection string
func (c *Config) BuildDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode, c.TimeZone,
	)
}

// NewDatabase creates a new database connection
func NewDatabase() (*gorm.DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Configure logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(config.BuildDSN()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Global database instance
var DB *gorm.DB

// Initialize sets up the database connection
func Initialize() error {
	db, err := NewDatabase()
	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}

	DB = db
	fmt.Println("Successfully connected to database")
	return nil
}

