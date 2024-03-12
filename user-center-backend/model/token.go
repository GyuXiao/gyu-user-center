package model

import (
	"context"
	"log"
	"user-center-backend/constant"
	"user-center-backend/pkg/errcode"
)

/*
* 使用 Redis 进行存放，校验，刷新，删除 token 信息
 */

var ctx = context.Background()

// 向 Redis 中插入一条 token 记录，记录 UserId 和 UserRole 信息

func InsertToken(token string, userId uint64, userRole int8) error {
	pipeline := RedisClient.TxPipeline()
	key := constant.TokenPrefixStr + token
	pipeline.HSet(ctx, key, constant.KeyUserId, userId, constant.KeyUserRole, userRole)
	pipeline.Expire(ctx, key, constant.TokenExpireTime)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Printf("redis insert token by userId err: %v", err)
		return err
	}
	return nil
}

// 判断 Redis 中是否存在对应的 token 记录

func CheckTokenExist(token string) ([]any, error) {
	key := constant.TokenPrefixStr + token
	result, err := RedisClient.HMGet(ctx, key, constant.KeyUserId, constant.KeyUserRole).Result()
	if err != nil {
		log.Printf("redis HMGet key err: %v", err)
		return nil, err
	}
	if result[0] == nil || result[1] == nil {
		return nil, errcode.ErrorTokenNotExist
	}
	return result, nil
}

// 刷新 token 的过期时间

func RefreshToken(token string) {
	_, err := CheckTokenExist(token)
	if err != nil {
		log.Printf("check redis key exist err: %v", err)
	}
	key := constant.TokenPrefixStr + token
	err = RedisClient.Expire(ctx, key, constant.TokenExpireTime).Err()
	if err != nil {
		log.Printf("redis expire key err: %v", err)
	}
}

// 删除 token

func DeleteToken(token string) error {
	key := constant.TokenPrefixStr + token
	return RedisClient.HDel(ctx, key, constant.KeyUserId, constant.KeyUserRole).Err()
}
