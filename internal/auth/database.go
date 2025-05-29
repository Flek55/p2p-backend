package auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func initDsn() string {
	return "host=" + os.Getenv("host") + "port=" + os.Getenv("port") + "user=" + os.Getenv("user") + "dbname=" + os.Getenv("dbname") + "password=" + os.Getenv("password") + "sslmode=disable"
}

func InitDb() *gorm.DB {
	godotenv.Load()
	dsn := initDsn()
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
