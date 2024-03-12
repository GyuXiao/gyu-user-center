package global

import (
	"user-center-backend/pkg/logger"
	"user-center-backend/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	RedisSetting    *setting.RedisSettingS
	Logger          *logger.Logger
	JWTSetting      *setting.JWTSettingS
	EmailSetting    *setting.EmailSettingS
)
