package route

import (
	"time"

	"go-server/api/middleware"
	"go-server/api/service"
	"go-server/domain"

	"go-server/bootstrap"

	"go-server/mongo"

	"go-server/repository"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")

	repository.InitCache()
	repository.DB = &db
	repository.ContextTimeout = timeout
	service.Env = env

	domain.InitJsons()

	publicRouter.POST("/signup", service.Signup)
	publicRouter.POST("/login", service.Login)
	publicRouter.POST("/refresh", service.RefreshToken)

	protectedRouter := gin.Group("")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	protectedRouter.POST("/autocoin", service.CoinAutoGrowing)
	protectedRouter.POST("/getofflinecoin", service.GetOfflineCoin)
	protectedRouter.POST("/clickearn", service.ClickEarn)
	//签到
	protectedRouter.POST("/checkin", service.CheckIn)
	//在线奖励
	protectedRouter.POST("/onlinerewards", service.ClaimOnlineRewards)
	protectedRouter.POST("/levelup", service.LevelUp)
	protectedRouter.POST("/unlockrole", service.UnLockRole)
	protectedRouter.POST("/passchapter", service.PassChapter)
	protectedRouter.POST("/ranking", service.Ranking)
	protectedRouter.POST("/caishen", service.CaiShen)
	protectedRouter.POST("/quickearn", service.QuickEarn)
	protectedRouter.POST("/continuousclick", service.ContinuousClick)
	protectedRouter.POST("/timesbonus", service.TimesBonus)
	//公寓小区
	protectedRouter.POST("/apartmentupgrade", service.UpgradeApartment)
	//坐骑
	protectedRouter.POST("/vehicleupgrade", service.UpgradeVehicle)
	protectedRouter.POST("/vehiclechange", service.ChangeVehicle)
	protectedRouter.POST("/unlockvehicle", service.UnLockVehicle)
	//资产
	protectedRouter.POST("/unlockcapital", service.UnLockCapital)
	protectedRouter.POST("/getcapitalincome", service.GetCapitalIncome)
	protectedRouter.POST("/sellcapital", service.SellCapital)

	protectedRouter.GET("/profile", service.Fetch)
	protectedRouter.GET("/task", service.CreateTask)
	protectedRouter.POST("/task", service.FetchTask)
}
