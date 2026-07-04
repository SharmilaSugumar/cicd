package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection using environment variables
func ConnectDatabase() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Default fallback for local development if env vars are not set
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "forgeflow"
	}
	if password == "" {
		password = "forgeflow"
	}
	if dbname == "" {
		dbname = "forgeflow"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established")
	return nil
}

// AutoMigrate runs GORM auto migrations for all models
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not established")
	}

	log.Println("Running auto-migration...")
	err := DB.AutoMigrate(
		&User{},
		&Organization{},
		&OrganizationMember{},
		&Project{},
		&Queue{},
		&Pipeline{},
		&PipelineRun{},
		&Job{},
		&JobDependency{},
		&JobLog{},
		&RetryPolicy{},
		&Worker{},
		&WorkerHeartbeat{},
		&DeadLetterQueue{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Auto-migration completed successfully")
	return nil
}
