package main

import (
	"github.com/gin-gonic/gin"
	"micro-book/microBook/internal/web"
)

func main() {
	server := gin.Default()
	u := &web.UserHandler{}
	u.RegisterRoutes(server)
	server.Run("localhost:8080")
}
