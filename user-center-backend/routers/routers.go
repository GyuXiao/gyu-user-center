package routers

import (
	"github.com/gin-gonic/gin"
	"time"
	"user-center-backend/global"
	v2 "user-center-backend/handlers/user/v2"
	"user-center-backend/middleware"
	"user-center-backend/pkg/limiter"
)

// TODO: 说明这里只对 /auth 接口进行限流
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimitBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(60 * time.Second))

	// 1,先刷新 Token
	r.Use(middleware.RefreshToken)

	user := v2.NewUser()
	apiv2 := r.Group("/api/user")
	// 普通用户
	// 用户注册
	apiv2.POST("/register", user.SignupHandler)
	// 用户登陆
	apiv2.POST("/login", user.LoginHandler)

	// 2,根据 Token 保存用户信息
	apiv2.Use(middleware.SetPersonalDetailsByToken)
	// 用户注销
	apiv2.POST("/logout", user.LogoutHandler)
	// 获取当前用户信息
	apiv2.GET("/current", user.CurrentUser)
	// 管理员
	// 查询用户
	apiv2.GET("/search", user.Search)
	// 删除用户
	apiv2.POST("/delete", user.Delete)

	return r
}
