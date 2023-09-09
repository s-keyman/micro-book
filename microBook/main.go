package main

import (
	"microBook/internal/web/middleware"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"microBook/internal/repository"
	"microBook/internal/repository/dao"
	"microBook/internal/service"
	"microBook/internal/web"
)

func main() {
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
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

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysid", store))
	// 校验步骤
	server.Use(middleware.NewLoginMiddleWare().IgnorePaths([]string{"/users/login", "/users/signup"}).Build())

	// v1
	//middleware.IgnorePaths = []string{"sss"}
	//server.Use(middleware.CheckLogin())
	////middleware.IgnorePaths = []string{"sss"}
	////server.Use(middleware.CheckLogin())
	//
	//// 不能忽略sss这条路径
	//server1 := gin.Default()
	//server1.Use(middleware.CheckLogin())
	////server1 := gin.Default()
	////server1.Use(middleware.CheckLogin())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/microBook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
