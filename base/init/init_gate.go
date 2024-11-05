package init

import (
	"fmt"
	config "iot_go/base/conf"
	sqlite "iot_go/base/db"
	logbase "iot_go/base/log"
	util "iot_go/base/util"
	"iot_go/server"
	"log"
	"os"
)

func Module(configPath string) error {
	return initConfig(configPath, []string{"base"})
}

func initConfig(configPath string, modules []string) error {
	if configPath == "" {
		fmt.Println("input config path is empty,you can put it like  ./dev/dev/")
		os.Exit(1)
	}
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")
	ips := util.GetLocalIPs()
	if len(ips) > 0 {
		config.LocalIP = ips[0]
	}
	if err := config.InitBaseConfig(configPath); err != nil {
		return err
	}
	if in := util.StrInArray("base", modules); in {
		logbase.BootstrapConf("base")
	}
	if util.StrInArray("sqlite3", modules) {
		if err := sqlite.InitSqlite3DB(config.GetConfigFilePath("sqlite3_map")); err != nil {
			fmt.Printf("[ERROR] InitDBPool:" + err.Error())
		}
	}
	server.HttpServerRun()
	return nil
}

func Destroy() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	server.HttpServerStop()
	logbase.Close()
	log.Printf("[INFO] %s\n", " success destroy resources.")
}
