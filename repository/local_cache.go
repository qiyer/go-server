package repository

import (
	"encoding/json"
	"time"

	"go-server/domain"

	"github.com/coocood/freecache"
)

var UserCache *freecache.Cache
var RankingCache *freecache.Cache
var LastLoginCache *freecache.Cache
var User_Cache_Time = 300      // Cache expiration time in seconds
var Last_Login_Cache_Time = 30 // Last login cache expiration time in seconds

func InitCache() {
	// Initialize the freecache with a size of 512MB
	UserCache = freecache.NewCache(512 * 1024 * 1024)
	RankingCache = freecache.NewCache(10 * 1024 * 1024)
	LastLoginCache = freecache.NewCache(100 * 1024 * 1024)
}

func GetUserCache(key string) (domain.User, error) {
	bytes, err := UserCache.Get([]byte(key))
	if err != nil {
		return domain.User{}, err // Return nil if the key does not exist
	}

	var newUser domain.User
	err = json.Unmarshal(bytes, &newUser) // 必须传递指针
	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func SetUserCache(key string, value domain.User) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	UserCache.Set([]byte(key), bytes, User_Cache_Time) // 0 means no expiration
}

func DeleteUserCache(key string) {
	UserCache.Del([]byte(key))
}

func GetLastLoginCache(key string) int64 {
	bytes, err := LastLoginCache.Get([]byte(key))
	if err != nil {
		return time.Now().Unix() // Return 0 if the key does not exist
	}

	var lastLogin int64
	err = json.Unmarshal(bytes, &lastLogin)
	if err != nil {
		return time.Now().Unix()
	}
	return lastLogin
}

func SetLastLoginCache(key string, value int64) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	LastLoginCache.Set([]byte(key), bytes, Last_Login_Cache_Time) // 0 means no expiration
}
