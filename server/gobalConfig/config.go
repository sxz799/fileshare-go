package gobalConfig

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	LimitFileLife   float64
	LimitFileSize   int
	ShareCodeType   int
	ShareCodeLength int
	ServerPort      string
	GinMode         string
)

func init() {
	log.Println("正在检查files文件夹")
	_, err := os.Stat("files")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir("files", os.ModePerm)
	}
}

func init() {
	log.Println("正在应用配置文件...")
	if _, err := os.Stat("conf/conf.yaml"); os.IsNotExist(err) {
		viper.SetConfigFile("conf.yaml")
	} else {
		viper.SetConfigFile("conf/conf.yaml")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("viper load fail ...")
		return
	}

	ServerPort = viper.GetString("server.port")
	GinMode = viper.GetString("server.ginMode")
	LimitFileLife = viper.GetFloat64("config.limitFileLife")
	LimitFileSize = viper.GetInt("config.limitFileSize")
	ShareCodeType = viper.GetInt("config.shareCodeType")
	ShareCodeLength = viper.GetInt("config.shareCodeLength")
}

func UseFrontMode(r *gin.Engine) {
	r.LoadHTMLGlob("dist/index.html")
	r.Static("/dist", "dist")
	r.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", "")
	})
}
