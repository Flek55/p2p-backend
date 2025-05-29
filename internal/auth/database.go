package auth

import (
	"fmt"
	"log"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func InitDb() *gorm.DB {
	dsn := "host=localhost port=5432 user=emilkerimov dbname=postgres password=123 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{}, &Session{})
	return db
}

func getDB() *gorm.DB {
	if database == nil {
		database = InitDb()
		for database == nil {
			database = InitDb()
		}
	}
	return database
}