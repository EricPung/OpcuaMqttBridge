package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open opens the database using sqlite.
func Open() (*gorm.DB, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "bridge.db"
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AutoMigrate migrates the models.
func AutoMigrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
