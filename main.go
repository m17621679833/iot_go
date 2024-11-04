package main

import (
	"flag"
	"fmt"
	conf "iot_go/base/conf"
	boot "iot_go/base/init"
	logger "iot_go/base/log"
	"time"
)

var (
	config = flag.String("config", "./env/dev/", "input config file like ./env/dev/")
)

func main() {
	flag.Parse()
	err := boot.Module(*config)
	if err != nil {
		return
	}
	logger.Warning("test message")
	logger.Close()
	time.Sleep(time.Second)
	fmt.Printf("  sssss%v\n", conf.Conf)
}
