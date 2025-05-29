package main

import (
	"context"
	"time"

	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	auth.InitDb()

	exUser := auth.User{
		ID:           uuid.New(),
		Email:        "popka",
		PasswordHash: "ggggg",
		CreatedAt:    time.Now(),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err := auth.CreateUser(ctx, &exUser); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Run()
}
