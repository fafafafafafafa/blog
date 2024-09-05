package main

import (
	"flag"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/middleware"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	//"os"
	//"log/slog"
)

func main() {
	// mkdirErr := os.MkdirAll("../public/uploaded", os.ModePerm) // 尝试创建存储路径
	// if mkdirErr != nil {
	//	slog.Error("function os.MkdirAll() Filed", slog.Any("err", mkdirErr.Error()))
	//return
	//}

	//out, createErr := os.Create("./hello.txt")
	//        if createErr != nil {
	//        slog.Error("function os.Create() Filed", slog.String("err", createErr.Error()))
	//        return
	//}
	//defer out.Close() // 创建文件 defer 关闭

	configPath := flag.String("c", "../config.yml", "配置文件路径")
	flag.Parse()

	// 根据命令行参数读取配置文件, 其他变量的初始化依赖于配置文件对象
	conf := g.ReadConfig(*configPath)

	_ = ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)

	// 初始化 gin 服务
	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})
	// 开发模式使用 gin 自带的日志和恢复中间件, 生产模式使用自定义的中间件
	if conf.Server.Mode == "debug" {
		r.Use(gin.Logger(), gin.Recovery()) // gin 自带的日志和恢复中间件, 挺好用的
	} else {
		r.Use(middleware.Recovery(true), middleware.Logger())
	}
	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithRedisDB(rdb))
	r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))
	ginblog.RegisterHandlers(r)

	serverAddr := conf.Server.Port
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://localhost:%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}
	r.Run(serverAddr)
}
