package server

import (
	"context"
	"github.com/gin-gonic/gin"
	config "iot_go/base/conf"
	"log"
	"net/http"
	"time"
)

var HttpSrvHandler *http.Server

func HttpServerRun() {
	gin.SetMode(config.Conf.BaseInfo.DebugMode)
	r := InitRoute()
	HttpSrvHandler = &http.Server{
		Addr:           config.Conf.HttpConfig.Addr,
		Handler:        r,
		ReadTimeout:    config.Conf.HttpConfig.ReadTimeout * time.Second,
		WriteTimeout:   config.Conf.HttpConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << uint(config.Conf.HttpConfig.MaxHeaderBytes),
	}
	go func() {
		log.Printf(" [INFO] HttpServerRun:%s\n", config.Conf.HttpConfig.Addr)
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", config.Conf.HttpConfig.Addr, err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 调用 Shutdown 方法进行优雅停机
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
