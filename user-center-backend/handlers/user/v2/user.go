package v2

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"user-center-backend/constant"
	"user-center-backend/global"
	"user-center-backend/model"
	"user-center-backend/pkg/app"
	"user-center-backend/pkg/errcode"
	"user-center-backend/pkg/util"
	"user-center-backend/service"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) SignupHandler(c *gin.Context) {
	// 参数校验
	param := service.UserSignupRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	// 业务处理——注册用户
	svc := service.New(c.Request.Context())
	Id, err := svc.Signup(&param)
	if err != nil {
		if errors.Is(err, errcode.ErrorUserRegisterParams) {
			response.ToErrorResponse(errcode.ErrorUserRegisterParams)
			return
		}
		if errors.Is(err, errcode.ErrorUserExit) {
			// 用户已经存在
			response.ToErrorResponse(errcode.ErrorUserExit)
			return
		}
		global.Logger.Errorf(c, "svc.Signup failed, err: %v", err)
		response.ToErrorResponse(errcode.ErrorUserSignupFail.WithDetails("用户注册时服务内部发生错误"))
		return
	}

	// 业务响应
	response.ToErrorResponse(errcode.Success.WithData(Id))
	return
}

func (u User) LoginHandler(c *gin.Context) {
	// 参数校验
	param := service.UserLoginRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	// 业务处理——用户登陆
	svc := service.New(c.Request.Context())
	user, err := svc.Login(&param)
	if err != nil {
		if errors.Is(err, errcode.ErrorUserLoginParams) {
			response.ToErrorResponse(errcode.ErrorUserLoginParams)
			return
		}
		if errors.Is(err, errcode.ErrorUserNotExit) {
			response.ToErrorResponse(errcode.ErrorUserNotExit)
			return
		}
		if errors.Is(err, errcode.ErrorUserPassword) {
			response.ToErrorResponse(errcode.ErrorUserPassword)
			return
		}
		global.Logger.Errorf(c, "svc.Login errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserLoginFail.WithDetails("用户登陆时服务内部错误"))
		return
	}

	// 记录用户登陆态
	session, err := json.Marshal(user)
	if err != nil {
		global.Logger.Errorf(c, constant.JSONMarshalError, err)
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("JSON Marshal 时发生错误"))
		return
	}
	c.SetCookie(constant.UserLoginState, string(session), constant.CookieExpire, "/", "", false, true)

	// 业务响应
	response.ToErrorResponse(errcode.Success.WithData(user))
	return
}

func (u User) LogoutHandler(c *gin.Context) {
	response := app.NewResponse(c)

	// 业务处理——用户注销
	// 讲 maxAge 的值设置为 -1 即为删除
	c.SetCookie(constant.UserLoginState, "", -1, "/", "", false, true)

	// 业务响应
	response.ToErrorResponse(errcode.Success.WithDetails("用户成功注销"))
	return
}

func (u User) CurrentUser(c *gin.Context) {
	response := app.NewResponse(c)
	// 从 cookie 从中拿到 user 对象
	session, err := c.Cookie(constant.UserLoginState)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("从 cookie 中获取信息时发生错误"))
		return
	}
	var user model.User
	// 将 session 反序列化为 user 对象
	err = json.Unmarshal([]byte(session), &user)
	if err != nil {
		global.Logger.Errorf(c, constant.JSONUnmarshalError, err)
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("JSON Unmarshal 时发生错误"))
		return
	}
	// 根据 user 对象中的字段，比如 userId，再查一遍数据库，然后脱敏返回
	latestUser, err := model.QueryUserByUserId(user.UserId)
	if err != nil {
		global.Logger.Errorf(c, "Query user by userId error: %v", err)
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("根据用户 id 查询最新用户信息时发生错误"))
		return
	}
	safeUser := util.GetSafetyUser(latestUser)
	response.ToErrorResponse(errcode.Success.WithData(safeUser))
	return
}

func (u User) Search(c *gin.Context) {
	response := app.NewResponse(c)

	// 从 cookie 从中拿到 user 对象
	session, err := c.Cookie(constant.UserLoginState)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("从 cookie 中获取信息时发生错误"))
		return
	}
	var user model.User
	// 将 session 反序列化为 user 对象
	err = json.Unmarshal([]byte(session), &user)
	if err != nil {
		global.Logger.Errorf(c, constant.JSONUnmarshalError, err)
		response.ToErrorResponse(errcode.ErrorCurrentUser.WithDetails("JSON Unmarshal 时发生错误"))
		return
	}

	// 根据用户名查询用户
	svc := service.New(c.Request.Context())
	users, err := svc.Search(user.Username)
	if err != nil {
		if errors.Is(err, errcode.ErrorSearchUserParams) {
			response.ToErrorResponse(errcode.ErrorSearchUserParams)
			return
		}
		if errors.Is(err, errcode.ErrorUserNotExit) {
			response.ToErrorResponse(errcode.ErrorUserNotExit)
			return
		}
		global.Logger.Errorf(c, "search user error: %v", err)
		response.ToErrorResponse(errcode.ErrorSearchUser.WithDetails("通过用户名查询用户时，服务内部发生错误"))
		return
	}

	// 业务响应
	response.ToErrorResponse(errcode.Success.WithData(users))
	return
}

func (u User) Delete(c *gin.Context) {
	response := app.NewResponse(c)

	// 根据 userId 进行删除用户（逻辑删除）
	svc := service.New(c.Request.Context())
	sid := c.Query("userId")
	uid, _ := strconv.ParseInt(sid, 10, 64)
	err := svc.Delete(uint64(uid))
	if err != nil {
		if errors.Is(err, errcode.ErrorUserNotExit) {
			response.ToErrorResponse(errcode.ErrorUserNotExit)
			return
		}
		global.Logger.Errorf(c, "delete user error: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteUser.WithDetails("删除用户时，服务内部发生错误"))
		return
	}

	// 业务响应
	response.ToErrorResponse(errcode.Success)
	return
}
