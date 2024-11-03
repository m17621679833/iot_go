package main

import (
	"flag"
	base "iot_go/base/conf"
	logger "iot_go/base/log"
	"time"
)

var (
	config = flag.String("config", "./env/dev/", "input config file like ./env/dev/")
)

func main() {
	flag.Parse()
	err := base.InitModule(*config)
	if err != nil {
		return
	}
	logger.Warning("test message")
	logger.Close()
	time.Sleep(time.Second)
}
