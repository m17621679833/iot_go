package base

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var ViperConfMap map[string]*viper.Viper
var Env string
var ConfPath string
var LocalIP = net.ParseIP("127.0.0.1")
var Conf *Base
var Sl3Conf *Sqlite3MapConf
var GORMMapPool map[string]*gorm.DB
var GORMDefaultPool *gorm.DB

type Base struct {
	BaseInfo  Setting   `mapstructure:"base"`
	LogConfig LogConfig `mapstructure:"log"`
}

type Setting struct {
	DebugMode    string `mapstructure:"debug_mode"`
	TimeLocation string `mapstructure:"time_location"`
}

type FileWriterConf struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type ConsoleWriterConf struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

type LogConfig struct {
	Level string            `mapstructure:"log_level"`
	FW    FileWriterConf    `mapstructure:"file_writer"`
	CW    ConsoleWriterConf `mapstructure:"console_writer"`
}

type Sqlite3MapConf struct {
	List map[string]*Sqlite3Conf `mapstructure:"list"`
}

type Sqlite3Conf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"db_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

func SetBaseConf(conf *Base) {
	Conf = conf
}

func SetSqlite3Conf(conf *Sqlite3MapConf) {
	Sl3Conf = conf
}

func SetConfEnvPath(confPath string) error {
	pathArr := strings.Split(confPath, "/")
	ConfPath = strings.Join(pathArr[:len(pathArr)-1], "/")
	Env = pathArr[len(pathArr)-2]
	return nil
}

func GetConfEnv() string {
	return Env
}

func GetConfigFilePath(fileName string) string {
	return ConfPath + "/" + fileName + ".toml"
}

func ParseConfigByFileName(fileName string, conf interface{}) error {
	file, err := os.Open(GetConfigFilePath(fileName))
	if err != nil {
		return fmt.Errorf("open config %v fail, %v", fileName, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail, %v", err)
	}
	v := viper.New()
	v.SetConfigType("toml")
	err = v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}

func ParseConfig(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %v fail, %v", path, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail, %v", err)
	}
	v := viper.New()
	v.SetConfigType("toml")
	err = v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}

// InitViperConf 初始化配置文件
func InitViperConf() error {
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

func InitBaseConfig(configPath string) error {
	if err := SetConfEnvPath(configPath); err != nil {
		return err
	}
	if err := InitViperConf(); err != nil {
		return err
	}
	return nil
}
