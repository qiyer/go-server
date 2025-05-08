package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(env *Env) *redis.Client {

	rHost := env.RedisHost
	rPort := env.RedisPort
	rDB := env.RedisDB
	rPass := env.RedisPass

	log.Println("rPass:", rPass)

	redisURI := fmt.Sprintf("%s:%s", rHost, rPort)
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       rDB,
	})

	// 检查 Redis 连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	return client
}
