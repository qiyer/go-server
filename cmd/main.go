package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	route "go-server/api/route"
	"go-server/bootstrap"
	"go-server/domain"
	"go-server/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	refreshTask()

	gin.Run(env.ServerAddress)
}

func refreshTask() {
	// 3. 创建 cron 实例（设置时区为东八区 UTC+8）
	// c := cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600))) // 北京时间
	c := cron.New(cron.WithSeconds()) // 支持秒级精度

	// 4. 添加每天零点执行的任务
	// "0 0 * * *" 表示 每天 00:00:00 执行
	_, err := c.AddFunc("40 11 17 * * *", dailyResetTask)
	// _, err := c.AddFunc("0 0 * * *", dailyResetTask)
	if err != nil {
		panic(err) // 如果 Cron 表达式错误，直接 panic
	}

	// 5. 启动 cron
	c.Start()
	fmt.Println("Cron 已启动，等待每天零点执行任务...")
	fmt.Println("执行每天 00:30 的任务:", time.Now().In(Location).Format("2006-01-02 15:04:05"))
}

var (
	// 加载 "Asia/Shanghai" 时区（北京时间 UTC+8）
	Location, _ = time.LoadLocation("Asia/Shanghai")
)

// 每日零点执行的任务
func dailyResetTask() {
	fmt.Println("执行每日零点刷新任务:", time.Now().Format("2006-01-02 15:04:05"))

	resp, err := http.Get("http://localhost:8080/ranking")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("读取响应体失败: %v", err)
		return
	}
	fmt.Printf("响应体:\n%s\n", string(body))
	var rankingResponse domain.RankingResponse
	if err := json.Unmarshal(body, &rankingResponse); err != nil {
		fmt.Errorf("解析 JSON 失败: %v", err)
		return
	}

	fmt.Println("排行榜数据：%+v\n", rankingResponse)

	repository.SetRankingCache("ranking", rankingResponse)
}
