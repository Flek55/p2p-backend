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
	return "host=" + os.Getenv("host") + " port=" + os.Getenv("port") + " user=" + os.Getenv("user") + " dbname=" + os.Getenv("dbname") + " password=" + os.Getenv("password") + " sslmode=disable"
}

func InitDb() {
	godotenv.Load()
	dsn := initDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{}, &Session{})
	database = db
}

func GetDB() *gorm.DB {
	if database == nil {
		for database == nil {
			InitDb()
		}
	}
	return database
}

func CreateUser(user *User) error {
	return GetDB().Create(user).Error
}
