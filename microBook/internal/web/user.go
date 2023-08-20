package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义所有跟用户有关的路由
type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	//注册路由
	//sever.POST("/users/signup", u.SignUp)
	//sever.POST("/users/login", u.Login)
	//sever.POST("/users/edit", u.Edit)
	//sever.POST("/users/profile", u.Profile)
	//使用分组功能
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	ctx.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(ctx *gin.Context)   {}
func (u *UserHandler) Edit(ctx *gin.Context)    {}
func (u *UserHandler) Profile(ctx *gin.Context) {}
