package main

import (
	"github.com/Flek55/p2p-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	auth.InitDb()

	r := gin.Default()

	r.Run()
}
