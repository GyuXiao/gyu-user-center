package middleware

import (
	"GyuBlog/constant"
	"GyuBlog/global"
	"GyuBlog/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 中间件鉴权

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie(constant.UserLoginState)
		if err == nil {
			if session != "" {
				var user model.User
				jsonErr := json.Unmarshal([]byte(session), &user)
				if jsonErr != nil {
					global.Logger.Errorf(c, "json Unmarshal error: %v", jsonErr)
					c.Abort()
					return
				}
				if user.UserRole != constant.RoleAdmin {
					global.Logger.Info("user is not administrator")
					c.Abort()
					return
				}
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		c.Abort()
	}
}
