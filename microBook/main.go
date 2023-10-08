package main

import (
	"microBook/internal/web/middleware"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
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
	err := server.Run(":8080")
	if err != nil {
		panic("系统故障")
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	//redisClient := redis.NewClient(
	//	&redis.Options{
	//		Addr:     "localhost:6379",
	//		Password: "", // no password set
	//		DB:       0,  // use default DB
	//	},
	//)
	// 基于滑动窗口算法的，利用redis进行IP限流
	//server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())
	server.Use(
		cors.New(
			cors.Config{
				//允许的单个路由（建议用 AllowOriginFunc ）
				//AllowOrigins:     []string{"https://foo.com"},
				//不填 AllowMethods 代表允许所有方法（ post put get 等，详见文档）
				AllowMethods: []string{"PUT", "POST", "GET", "OPTIONS"},
				AllowHeaders: []string{"Content-Type", "Authorization"},
				//是否允许带 cookie 之类的东西，不加前端拿不到 token
				ExposeHeaders:    []string{"x-jwt-token"},
				AllowCredentials: true,
				AllowOriginFunc: func(origin string) bool {
					if strings.HasPrefix(origin, "http://192.168.31.37") {
						//开发环境
						return true
					}
					return origin == "http://locolhost"
				},
				MaxAge: 12 * time.Second,
			},
		),
	)

	// 步骤1
	//store := cookie.NewStore([]byte("secret"))
	//store := cookie.NewStore([]byte("secret"))
	//使用 memstore 存储 session
	store := memstore.NewStore(
		[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"),
	)
	// 使用 redis 存储 session
	//store, err := redis.NewStore(
	//	16, "tcp", "localhost:6379", "",
	//	[]byte("eW*ZAxyp1Lx81hp9:swB?Sp)l$We8qeI"), []byte("QvcBUP5f[DTp!u>4G%x?atz@1d/}!DS^"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("mysid", store))
	// 校验步骤
	//server.Use(middleware.NewLoginMiddleWare().IgnorePaths([]string{"/users/login", "/users/signup"}).Build())
	// 使用 jwt
	server.Use(middleware.NewLoginJWTMiddleWare().IgnorePaths([]string{"/users/login", "/users/signup"}).Build())
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
