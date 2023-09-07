package web

import (
	"errors"
	"microBook/internal/domain"
	"microBook/internal/service"
	"net/http"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicate = service.ErrDuplicateEmail
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler 定义所有跟用户有关的路由
type UserHandler struct {
	svc *service.UserService
	//预编译正则表达式，不用每次使用前都编译一次
	emailRegeExp    *regexp2.Regexp
	passwordRegeExp *regexp2.Regexp
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

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:             svc,
		emailRegeExp:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegeExp: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
	}
}

// SignUp 用户注册接口
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析数据
	// 解析错误，直接写回一个 4xx 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//校验邮箱
	ok, err := u.emailRegeExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "输入的邮箱格式不对！")
		return
	}

	//校验密码
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致！")
		return
	}
	ok, err = u.passwordRegeExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含字母、数字、特殊字符")
		return
	}

	// 调用一下 service 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, ErrUserDuplicate) {
		ctx.String(http.StatusOK, "邮箱冲突！")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常！")
		return
	}

	ctx.String(http.StatusOK, "注册成功！")
}
func (u *UserHandler) Login(ctx *gin.Context) {
	//1.登录本身
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//2.登录态的校验
}
func (u *UserHandler) Edit(ctx *gin.Context)    {}
func (u *UserHandler) Profile(ctx *gin.Context) {}
