package main

import (
	"fmt"
	_ "go-gin-example/docs"
	"go-gin-example/models"
	"go-gin-example/pkg/gredis"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
	"log"
	"net/http"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin

// @license.name MIT

func main() {
	// endless.DefaultReadTimeOut = setting.ReadTimeout
	// endless.DefaultWriteTimeOut = setting.WriteTimeout
	// endless.DefaultMaxHeaderBytes = 1 << 20
	// endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
	// server := endless.NewServer(endPoint, routers.InitRouter())

	// server.BeforeBegin = func(add string) {
	// 	// logging.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }

	// 初始化
	setting.Setup()
	models.Setup()
	gredis.Setup()
	logging.Setup()

	r := routers.InitRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        r,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
