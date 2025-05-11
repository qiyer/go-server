package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/bootstrap"
	"go-server/domain"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 缓存用户数据
func CacheUserData(user *domain.User) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", user.ID)
	jsonData, _ := json.Marshal(user)

	log.Println("CacheUserData key:", key)

	return bootstrap.RedisClient.Set(ctx, key, jsonData, 30*time.Minute).Err()
}

// 获取缓存用户数据
func GetUserFromCache(userID primitive.ObjectID) (*domain.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", userID)
	log.Println("GetUserFromCache key:", key)
	data, err := bootstrap.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// 缓存用户数据
func CacheUserDataByOpenID(user *domain.User, openId string) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", openId)
	jsonData, _ := json.Marshal(user)
	return bootstrap.RedisClient.Set(ctx, key, jsonData, 30*time.Minute).Err()
}

// 获取缓存用户数据
func GetUserFromCacheByOpenID(openId string) (*domain.User, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:data:%s", openId)
	data, err := bootstrap.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
