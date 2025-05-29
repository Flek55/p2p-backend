package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	auth.InitDb()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s := auth.CreateService("hui", 300)

	err := s.Register(ctx, "burito1533@gmail.com", "zhara")

	if err != nil {
		fmt.Println("error creating user: ", err)
	}

	r := gin.Default()
	r.Run()
}
