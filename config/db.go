package config

import (
	// Go imports
	"fmt"
	"log"

	// External imports
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbConfig Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBSSLMode,
	)), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
