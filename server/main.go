package main

import (
	"fileshare-server/gobalConfig"
	"fileshare-server/model"
	"fileshare-server/router"
	"fileshare-server/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"log"
	"os"
)

func main() {
	log.Println("正在连接数据库...")
	util.InitDB()
	log.Println("正在检查表结构...")
	model.InitAutoMigrateDB()
	gin.SetMode(gobalConfig.GinMode)
	r := gin.Default()
	_, err := os.Stat("dist")
	if err == nil {
		gobalConfig.UseFrontMode(r)
		log.Println("已开启前后端整合模式！")
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowFiles = true

	r.Use(cors.New(config))

	router.RegRouter(r)

	c := cron.New()
	c.AddFunc("@every 10m", model.AutoDelFile)
	c.Start()
	log.Println("定时任务启动成功,服务启动成功,当前使用端口：", gobalConfig.ServerPort)
	err = r.RunTLS(":"+gobalConfig.ServerPort, "cert/pem.pem", "cert/key.key")
	if err != nil {
		log.Println("未找到证书文件，以http运行！")
		r.Run(":" + gobalConfig.ServerPort)
	}
}
