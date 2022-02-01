package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"os"
	"time"
)

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("kontak.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error gorm open: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// log
	if os.Getenv("ENV") != "PROD" {
		db = db.Debug()
	}
	return db, nil
}
