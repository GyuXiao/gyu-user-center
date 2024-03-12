package constant

import "time"

const (
	RoleUser  = 0
	RoleAdmin = 1

	UserLoginState = "user_session"
	CookieExpire   = 86400
	PatternStr     = "/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/"
)

const (
	JSONMarshalError   = "JSON marshal error: %v"
	JSONUnmarshalError = "JSON Unmarshal error: %v"
)

const (
	JwtSecret = "user-center-backend"
)

const (
	KeyUserId       = "user_id"
	KeyUserRole     = "user_role"
	TokenExpireTime = time.Hour * 24
	TokenPrefixStr  = "login:token:"
	TokenHeader     = "Authorization"
	TokenEmpty      = ""
)
