package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"micro-book/microBook/internal/web"
)

func main() {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//允许的单个路由（建议用 AllowOriginFunc ）
		//AllowOrigins:     []string{"https://foo.com"},
		//不填 AllowMethods 代表允许所有方法（ post put get 等，详见文档）
		AllowMethods: []string{"PUT", "POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//是否允许带 cookie 之类的东西
		//ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//开发环境
				return true
			}
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Second,
	}))
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
