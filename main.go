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
)

func main() {
	log.Println("正在连接数据库...")
	util.InitDB()
	log.Println("正在检查表结构...")
	model.InitAutoMigrateDB()
	gin.SetMode(gobalConfig.GinMode)
	r := gin.Default()
	if gobalConfig.FrontMode {
		log.Println("已开启前后端整合模式！")
		gobalConfig.UseFrontMode(r)
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	router.RegRouter(r)

	c := cron.New()
	c.AddFunc("@every 10m", model.AutoDelFile)
	c.Start()
	log.Println("定时任务启动成功,服务启动成功,当前使用端口：", gobalConfig.ServerPort)
	err := r.RunTLS(":"+gobalConfig.ServerPort, "cert/pem.pem", "cert/key.key")
	if err != nil {
		log.Println(err)
	}
}
