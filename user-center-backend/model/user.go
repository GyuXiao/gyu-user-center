package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"user-center-backend/global"
	"user-center-backend/pkg/errcode"
)

// MySQL 存放 User 数据信息

type User struct {
	// id
	Id int64 `gorm:"column:id"`
	// 用户 ID
	UserId uint64 `json:"userId" gorm:"column:userID" `
	// 用户昵称
	Username string `json:"username" gorm:"column:username" `
	// 登陆账号
	UserAccount string `json:"userAccount" gorm:"column:userAccount" `
	// 用户头像
	AvatarUrl string `json:"avatarUrl" gorm:"column:avatarUrl" `
	// 性别
	Gender int64 `json:"gender" gorm:"column:gender" `
	// 用户密码
	UserPassword string `json:"userPassword" gorm:"column:userPassword" `
	// 电话
	Phone string `json:"phone" gorm:"column:phone" `
	// 邮箱
	Email string `json:"email" gorm:"column:email" `
	// 用户状态 0-正常
	UserStatus int64 `json:"userStatus" gorm:"column:userStatus" `
	// 角色
	UserRole int8 `json:"userRole" gorm:"column:userRole"`
	// 默认字段
	CreateTime global.JsonTime `gorm:"column:createTime;autoCreateTime"`
	UpdateTime global.JsonTime `gorm:"column:updateTime;autoUpdateTime"`
	IsDelete   int8            `json:"isDelete" gorm:"column:isDelete"`
}

type UserFrontObject struct {
	UserId      uint64 `json:"userId,string"`
	Username    string `json:"username"`
	UserAccount string `json:"userAccount"`
	Token       string `json:"token,omitempty"`
	UserRole    int8   `json:"userRole"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Gender      int64  `json:"gender"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Create() (int64, error) {
	err := DBEngine.Model(&User{}).Create(&u).Error
	if err != nil {
		return -1, err
	}
	return u.Id, nil
}

func CheckUserExist(userAccount string) error {
	_, err := QueryUserByAccount(userAccount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return errcode.ErrorUserExit
}

func QueryUserByAccount(userAccount string) (*User, error) {
	var user User
	err := DBEngine.Where("isDelete=0 and userAccount = ?", userAccount).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func QueryUserByUsername(username string) (*[]User, error) {
	var user []User
	// 如果 username 是索引列的话，则不应该做左模糊查询（导致索引失效，全表扫描）
	err := DBEngine.Where("isDelete=0 and username like ?", "%"+username+"%").Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func QueryUserByUserId(userId uint64) (*User, error) {
	var user User
	err := DBEngine.Where("isDelete=0 and userId=?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByUserId(userId uint64) error {
	_, err := QueryUserByUserId(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errcode.ErrorUserNotExit
	}
	err = DBEngine.Model(&User{}).Where("isDelete=0 and userId=?", userId).Update("isDelete", 1).Error // 这里一定要用 Model，不然找不到 user 表
	if err != nil {
		return err
	}
	return nil
}
