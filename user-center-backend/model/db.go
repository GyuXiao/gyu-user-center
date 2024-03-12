package model

import (
	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
)

var DBEngine *gorm.DB

var RedisClient *redis.Client
