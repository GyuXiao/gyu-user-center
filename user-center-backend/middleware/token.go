package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"user-center-backend/constant"
	"user-center-backend/model"
	"user-center-backend/pkg/app"
	"user-center-backend/pkg/errcode"
)

var token = ""

// RefreshToken 刷新 token 的过期时间
func RefreshToken(c *gin.Context) {
	// 如果 header 没有 token 的话，跳过
	// 如果 header 有 token 的话，调用 redis 的 RefreshToken 函数
	token = c.GetHeader(constant.TokenHeader)
	if token == constant.TokenEmpty {
		c.Next()
		return
	}
	model.RefreshToken(token)
}

// SetPersonalDetailsByToken
// 通过 token 从 redis 拿取 user_id 和 user_role
// 用户登陆成功后，才会调用该函数
func SetPersonalDetailsByToken(c *gin.Context) {
	response := app.NewResponse(c)
	// 先判断 token 是否存在
	result, err := model.CheckTokenExist(token)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorTokenFetch)
		c.Abort()
		return
	}
	// 如果存在的话，从 redis 拿到 user_id 和 user_role 并保存
	id, ok := result[0].(string)
	if !ok {
		log.Println("valid user_id failed")
		response.ToErrorResponse(errcode.ErrorUserLoginFail)
		c.Abort()
		return
	}
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("parse user_id to int64 failed")
		response.ToErrorResponse(errcode.ErrorUserLoginFail)
		c.Abort()
		return
	}

	role, ok := result[1].(string)
	if !ok {
		log.Println("valid user_role failed")
		response.ToErrorResponse(errcode.ErrorUserLoginFail)
		c.Abort()
		return
	}
	userRole, err := strconv.Atoi(role)
	if err != nil {
		log.Println("parse user_role to int failed")
		response.ToErrorResponse(errcode.ErrorUserLoginFail)
		c.Abort()
		return
	}

	c.Set(constant.KeyUserId, userId)
	c.Set(constant.KeyUserRole, userRole)
	c.Next()
}
