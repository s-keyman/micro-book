package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddleWare() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(paths []string) *LoginMiddlewareBuilder {
	l.paths = paths
	//for _, path := range paths {
	//	l.paths = append(l.paths, path)
	//}
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 不需要登录校验的
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime, _ := ctx.Get("update_time")
		sess.Set("userId", id)
		sess.Options(
			sessions.Options{
				// 60 秒过期
				MaxAge: 1800,
			},
		)
		now := time.Now().UnixMilli()
		// 说明还没有刷新过，刚登陆，还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				panic(err)
			}
		}
		// updateTime 是有的，证明登录已经超过 1 分钟，需要刷新 session
		updateTimeVal, _ := updateTime.(int64)
		if now-updateTimeVal > 60000 {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				panic(err)
			}
		}
	}
}
