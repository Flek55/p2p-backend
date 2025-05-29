package auth

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
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

func CreateUser(ctx context.Context, user *User) error {
	return GetDB().WithContext(ctx).Create(user).Error
}

func FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func CreateSession(ctx context.Context, session *Session) error {
	return GetDB().WithContext(ctx).Create(session).Error
}

func DeleteSession(ctx context.Context, session *Session) error {
	return GetDB().WithContext(ctx).Delete(session).Error
}

func FindSessionByID(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	var session Session
	err := GetDB().WithContext(ctx).First(&session, "id = ?", sessionID).Error
	return &session, err
}
