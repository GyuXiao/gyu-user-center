package util

import (
	"GyuBlog/global"
	"GyuBlog/model"
)

func GetSafetyUser(user *model.User) model.User {
	if user == nil {
		global.Logger.Error("此时用户数据为空")
		return model.User{}
	}
	safetyUser := model.User{}
	safetyUser.Id = user.Id
	safetyUser.UserId = user.UserId
	safetyUser.Username = user.Username
	safetyUser.UserAccount = user.UserAccount
	safetyUser.AvatarUrl = user.AvatarUrl
	safetyUser.Gender = user.Gender
	safetyUser.Phone = user.Phone
	safetyUser.Email = user.Email
	safetyUser.UserStatus = user.UserStatus
	safetyUser.UserRole = user.UserRole
	safetyUser.CreateTime = user.CreateTime
	return safetyUser
}
