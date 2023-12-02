package routers

import (
	"github.com/gin-gonic/gin"
	"time"
	"user-center-backend/global"
	v2 "user-center-backend/handlers/user/v2"
	"user-center-backend/middleware"
	"user-center-backend/pkg/limiter"
)

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

	user := v2.NewUser()

	//r.GET("/auth", api.GetAuth)
	apiv2 := r.Group("/api/user")

	// 用户模块
	// 用户注册
	apiv2.POST("/register", user.SignupHandler)
	// 用户登陆
	apiv2.POST("/login", user.LoginHandler)
	// 用户注销
	apiv2.POST("/logout", user.LogoutHandler)
	// 获取当前用户信息
	apiv2.GET("/current", user.CurrentUser)

	// 管理员
	// 查询用户
	apiv2.GET("/search", middleware.AuthMiddleWare(), user.Search)
	// 删除用户
	apiv2.POST("/delete", middleware.AuthMiddleWare(), user.Delete)

	//apiv2.Use(app.JWT())
	return r
}
