package web

import (
	"errors"
	"microBook/internal/domain"
	"microBook/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
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
	err = u.svc.SignUp(
		ctx, domain.User{
			Email:    req.Email,
			Password: req.Password,
		},
	)
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
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	// Bind 方法会根据 Content-Type 来解析数据
	// 解析错误，直接写回一个 4xx 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}
	// 步骤2
	// 在这里登录成功了
	// 设置 session
	ssid := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	ssid.Set("userId", user.Id)
	ssid.Options(
		sessions.Options{
			//设置过期时间
			MaxAge: 60,
		},
	)
	err = ssid.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

// LoginJWT jwt登录
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type jwtReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req jwtReq
	// Bind 方法会根据 Content-Type 来解析数据
	// 解析错误，直接写回一个 4xx 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}

	//步骤2 jwt 设置登录态
	//生成一个 JWT token
	token := jwt.New(jwt.SigningMethodHS512)
	tokenString, err := token.SignedString([]byte("eW*ZAxyp1Lx81hp9:swB?Sp)l$We8qeI"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}
	ctx.Header("x-jwt-token", tokenString)

	ctx.String(http.StatusOK, strconv.FormatUint(user.Id, 10)+"\n")

	ctx.String(http.StatusOK, "登录成功")
}

// Edit 用户编译信息
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		// 注意，其它字段，尤其是密码、邮箱和手机，
		// 修改都要通过别的手段
		// 邮箱和手机都要验证
		// 密码更加不用多说了
		Nickname string `json:"nickname"`
		// 2023-01-01
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req EditReq
	// Bind 方法会根据 Content-Type 来解析数据
	// 解析错误，直接写回一个 4xx 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 你可以尝试在这里校验。
	// 比如说你可以要求 Nickname 必须不为空
	// 校验规则取决于产品经理
	if req.Nickname == "" {
		ctx.String(http.StatusOK, "昵称不能为空")
		return
	}

	// 限制个人简介的长度
	if len(req.AboutMe) > 1024 {
		ctx.String(http.StatusOK, "个人简介不能超过1024个字符！")
		return
	}

	//日期转换
	birthday, err := time.Parse(time.DateTime, req.Birthday)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.String(http.StatusOK, "日期格式不对")
		return
	}

	err = u.svc.UpdateNonSensitiveInfo(
		ctx, domain.User{
			Nickname: req.Nickname,
			AboutMe:  req.AboutMe,
			Birthday: birthday},
	)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "OK")
	return
}

// Profile 用户详情
func (u *UserHandler) Profile(ctx *gin.Context) {
	type ProfileReq struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}

	var req ProfileReq
	// Bind 方法会根据 Content-Type 来解析数据
	// 解析错误，直接写回一个 4xx 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ctx.String(http.StatusOK, "成功")
	return
}
