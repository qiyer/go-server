package route

import (
	"time"

	"go-server/api/middleware"
	"go-server/api/service"

	"go-server/bootstrap"

	"go-server/mongo"

	"go-server/repository"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")

	repository.DB = &db
	repository.ContextTimeout = timeout
	service.Env = env

	publicRouter.POST("/signup", service.Signup)
	publicRouter.POST("/login", service.Login)
	publicRouter.POST("/refresh", service.RefreshToken)

	protectedRouter := gin.Group("")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	protectedRouter.POST("/autocoin", service.CoinAutoGrowing)
	protectedRouter.POST("/checkin", service.CheckIn)
	protectedRouter.POST("/onlinerewards", service.ClaimOnlineRewards)
	protectedRouter.POST("/levelup", service.LevelUp)
	protectedRouter.POST("/passchapter", service.PassChapter)
	protectedRouter.POST("/ranking", service.Ranking)

	protectedRouter.GET("/profile", service.Fetch)
	protectedRouter.GET("/task", service.CreateTask)
	protectedRouter.POST("/task", service.FetchTask)
}
