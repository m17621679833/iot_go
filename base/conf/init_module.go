package base

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	logFaced "iot_go/base/log"
	sqlite3 "iot_go/base/sqlite3"
	trace "iot_go/base/trace"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var TimeLocation *time.Location
var TimeFormat = "2006-01-02 15:04:05"
var DateFormat = "2006-01-02"
var LocalIP = net.ParseIP("127.0.0.1")

func InitModule(configPath string) error {
	return initModule(configPath, []string{"base"})
}

func initModule(configPath string, modules []string) error {
	if configPath == "" {
		fmt.Println("input config file like ./conf/dev/")
		os.Exit(1)
	}
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIP = ips[0]
	}
	ParseConfEnvPath(configPath)
	if err := InitViperMap(); err != nil {
		return err
	}
	if InArrayString("base", modules) {
		if err := InitBaseConf(GetConfPath("base")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}
	if InArrayString("sqlite3", modules) {
		if err := sqlite3.InitSqlite3DB(GetConfPath("sqlite3_map")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitDBPool:"+err.Error())
		}
	}
	if location, err := time.LoadLocation(ConfBase.Base.TimeLocation); err != nil {
		return err
	} else {
		TimeLocation = location
	}
	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}

func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

func InitBaseConf(path string) error {
	ConfBase = &Conf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
		return err
	}
	if ConfBase.Base.DebugMode == "" {
		if ConfBase.Base.DebugMode != "" {
			ConfBase.Base.DebugMode = ConfBase.Base.DebugMode
		} else {
			ConfBase.Base.DebugMode = "debug"
		}
	}
	if ConfBase.Base.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.Base.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.Base.TimeLocation = "Asia/Chongqing"
		}
	}
	if ConfBase.Log.Level == "" {
		ConfBase.Log.Level = "trace"
	}
	//配置日志
	logConf := logFaced.LogConfig{
		Level: ConfBase.Log.Level,
		FW: logFaced.FileWriterConf{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
		},
		CW: logFaced.ConsoleWriterConf{
			On:    ConfBase.Log.CW.On,
			Color: ConfBase.Log.CW.Color,
		},
	}
	err = logFaced.InitLogConf(logConf)
	if err != nil {
		return err
	}
	logFaced.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func GetSqliteDBWithTrace(c *trace.TraceContext) *gorm.DB {
	db := sqlite3.GORMDefaultPool
	db = db.Set("trace_context", c)
	return db
}

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}

// InitViperMap 初始化配置文件
func InitViperMap() error {
	f, err := os.Open(ConfPath + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}
