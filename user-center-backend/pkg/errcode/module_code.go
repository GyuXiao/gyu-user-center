// Package errcode
/**
  @author: zk.xiao
  @date: 2022/6/1
  @note:
**/
package errcode

// User 错误码

var (
	ErrorUserSignupFail     = NewError(20010001, "用户注册失败")
	ErrorUserExit           = NewError(20010002, "用户已经存在")
	ErrorUserNotExit        = NewError(20010003, "用户不存在")
	ErrorUserPassword       = NewError(20010004, "用户密码错误")
	ErrorUserLoginFail      = NewError(20010005, "用户登陆失败")
	ErrorUserRegisterParams = NewError(20010006, "用户注册参数错误")
	ErrorUserLoginParams    = NewError(20010007, "用户登陆参数错误")
	ErrorCurrentUser        = NewError(20010008, "获取当前用户信息错误")
	ErrorUserNoLogin        = NewError(20010009, "用户未登录")
	ErrorUserLogout         = NewError(20010010, "用户退出失败")
	ErrorUserIdNotExist     = NewError(20010011, "用户 ID 不存在")
	ErrorUserIdConvert      = NewError(20010012, "用户 ID 格式转换错误")
	ErrorUserRoleNotExist   = NewError(20010013, "用户角色不存在")
	ErrorUserRoleConvert    = NewError(20010014, "用户角色格式转换错误")
)

// Token 错误码

var (
	ErrorTokenNotExist = NewError(30010001, "Token 不存在")
	ErrorTokenFetch    = NewError(30010002, "从 Redis 获取 Token 失败")
)

// 管理员 错误码

var (
	ErrorSearchUser       = NewError(40010001, "搜索用户错误")
	ErrorSearchUserParams = NewError(40010002, "搜索用户参数错误")
	ErrorDeleteUser       = NewError(40010003, "删除用户错误")
)
