package middleware

import (
	"log"
	"microBook/internal/web"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddleWare() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(paths []string) *LoginJWTMiddlewareBuilder {
	l.paths = paths
	//for _, path := range paths {
	//	l.paths = append(l.paths, path)
	//}
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		// 没有登录
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			//有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 里面一定要传指针
		token, err := jwt.ParseWithClaims(
			tokenStr, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("eW*ZAxyp1Lx81hp9:swB?Sp)l$We8qeI"), nil
			},
		)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 每 10 秒续约一次，生成了一个新的token
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenString, err := token.SignedString([]byte("eW*ZAxyp1Lx81hp9:swB?Sp)l$We8qeI"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenString)
		}
		ctx.Set("claims", claims)
	}
}
