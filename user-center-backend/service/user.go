package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"regexp"
	"time"
	"user-center-backend/global"
	"user-center-backend/model"
	"user-center-backend/pkg/errcode"
	"user-center-backend/pkg/snowflake"
	"user-center-backend/pkg/util"
)

func (svc *Service) Signup(p *UserSignupRequest) (int64, error) {
	// 业务校验参数
	if p.UserAccount == "" || p.UserPassword == "" || p.ConfirmPassword == "" {
		return -1, errcode.ErrorUserRegisterParams
	}
	if len(p.UserAccount) < 6 || len(p.UserPassword) < 8 || len(p.ConfirmPassword) < 8 {
		return -1, errcode.ErrorUserRegisterParams
	}
	_, err := regexp.MatchString("/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/", p.UserAccount)
	if err != nil {
		return -1, errcode.ErrorUserRegisterParams
	}
	// 账号不能重复
	// 先判断待注册的用户的用户名是否已经存在
	err = model.CheckUserExist(p.UserAccount)
	if err != nil || errors.Is(err, errcode.ErrorUserExit) {
		return -1, err
	}

	// 通过雪花算法获取 userID
	userID, snowErr := snowflake.GetID()
	if snowErr != nil {
		return -1, snowErr
	}
	// 先对密码加密
	pwd, err := encodePassword(p.UserPassword)
	if err != nil {
		return -1, err
	}

	u := &model.User{
		UserId:       userID,
		UserAccount:  p.UserAccount,
		UserPassword: pwd,
		CreateTime:   global.JsonTime(time.Now()),
		UpdateTime:   global.JsonTime(time.Now()),
	}
	//注册用户
	return u.Create()
}

func encodePassword(pwd string) (string, error) {
	hashStr, err := util.EncodeBcrypt(pwd)
	if err != nil {
		return "", err
	}
	return util.EncodeMd5([]byte(hashStr)), nil
}

func (svc *Service) Login(p *UserLoginRequest) (user *model.User, error error) {
	// 1，参数校验
	if p.UserAccount == "" || p.Password == "" {
		return nil, errcode.ErrorUserLoginParams
	}
	if len(p.UserAccount) < 6 || len(p.Password) < 8 {
		return nil, errcode.ErrorUserLoginParams
	}
	_, err := regexp.MatchString("/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/", p.UserAccount)
	if err != nil {
		return nil, errcode.ErrorUserLoginParams
	}

	// 2，校验密码是否正确（数据库查询并比对）
	u, err := model.QueryUserByAccount(p.UserAccount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrorUserNotExit
		}
		return nil, err
	}
	pwd := util.DecodeMd5(u.UserPassword)
	if !util.DecodeBcrypt(pwd, p.Password) {
		return nil, errcode.ErrorUserPassword
	}

	// 3，数据查询成功后，返回脱敏后的用户信息
	safeUser := util.GetSafetyUser(u)
	return &safeUser, nil

	// 返回  accessToken 和 refreshToken
	/*
		// 数据查询成功
		user = &model.User{
			UserName: p.UserAccount,
			UserID:   userID,
		}
		accessToken, refreshToken, genError := jwt.GenToken(user.UserID, user.UserName)
		if genError != nil {
			return nil, genError
		}
		user.AccessToken = accessToken
		user.RefreshToken = refreshToken
		return
	*/

}

func (svc *Service) Search(username string) ([]model.User, error) {
	// 参数校验
	if len(username) < 6 {
		return nil, errcode.ErrorSearchUserParams
	}
	_, err := regexp.MatchString("/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/", username)
	if err != nil {
		return nil, errcode.ErrorSearchUserParams
	}

	// 数据库查询
	users, err := model.QueryUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrorUserNotExit
		}
		return nil, err
	}

	// 返回脱敏信息
	var safeUsers []model.User
	for _, user := range *users {
		safeUser := util.GetSafetyUser(&user)
		safeUsers = append(safeUsers, safeUser)
	}
	return safeUsers, nil
}

func (svc *Service) Delete(uid uint64) error {
	err := model.DeleteUserByUserId(uid)
	if err != nil {
		return err
	}
	return nil
}
