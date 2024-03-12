package util

import (
	"github.com/gin-gonic/gin"
	"user-center-backend/constant"
	"user-center-backend/pkg/errcode"
)

func GetUserId(c *gin.Context) (int64, error) {
	value, exist := c.Get(constant.KeyUserId)
	if !exist {
		return -1, errcode.ErrorUserIdNotExist
	}
	userId, ok := value.(int64)
	if !ok {
		return -1, errcode.ErrorUserIdConvert
	}
	return userId, nil
}

func GetUserRole(c *gin.Context) (int, error) {
	value, exist := c.Get(constant.KeyUserRole)
	if !exist {
		return -1, errcode.ErrorUserRoleNotExist
	}
	userRole, ok := value.(int)
	if !ok {
		return -1, errcode.ErrorUserRoleConvert
	}
	return userRole, nil
}
