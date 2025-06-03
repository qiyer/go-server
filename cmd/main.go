package main

import (
	"time"

	route "go-server/api/route"
	"go-server/bootstrap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	// 配置 CORS 中间件
	gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Custom-Header"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.Setup(env, timeout, db, gin)

	gin.Run(env.ServerAddress)
}
