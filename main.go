package main

import (
	"flag"
	boot "iot_go/base/init"
	"os"
	"os/signal"
	"syscall"
)

var (
	config = flag.String("config", "./env/dev/", "input config file like ./env/dev/")
)

func main() {
	flag.Parse()
	err := boot.Goo(*config)
	if err != nil {
		return
	}
	defer boot.Destroy()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
