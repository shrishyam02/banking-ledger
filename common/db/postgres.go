package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres establishes a connection to the PostgreSQL database.
func ConnectPostgres(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Check connection
	if err := db.Raw("SELECT 1").Error; err != nil {
		return nil, err
	}
	log.Println("Successfully connected to Postgres")
	return db, nil
}
