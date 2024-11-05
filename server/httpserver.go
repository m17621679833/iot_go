package server

import (
	"context"
	"github.com/gin-gonic/gin"
	config "iot_go/base/conf"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(config.GetStringConf("env.base.debug_mode"))
	r := InitRoute()
	HttpSrvHandler := &http.Server{
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
	/*go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C: // 当定时器触发时
				printMemStats() // 打印内存使用统计信息
			}
		}
	}()*/
}

/*
	func printMemStats() {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		allocMB := float64(memStats.Alloc) / (1024 * 1024)
		totalAllocMB := float64(memStats.TotalAlloc) / (1024 * 1024)
		sysMB := float64(memStats.Sys) / (1024 * 1024)
		numGC := memStats.NumGC
		log.Printf("Alloc:%.2f MB,TotalAlloc:%.2f MB,Sys:%.2f MB,NumGC:%d", allocMB, totalAllocMB, sysMB, numGC)
	}
*/
func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
