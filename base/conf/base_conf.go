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
	"time"
)

var ViperConfMap map[string]*viper.Viper
var Env string
var ConfPath string
var LocalIP = net.ParseIP("127.0.0.1")
var Conf *Base
var GORMMapPool map[string]*gorm.DB
var GORMDefaultPool *gorm.DB

type Base struct {
	BaseInfo      Setting       `mapstructure:"base"`
	LogConfig     LogConfig     `mapstructure:"log"`
	HttpConfig    HttpConfig    `mapstructure:"http"`
	SwaggerConfig SwaggerConfig `mapstructure:"swagger"`
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
	LogLevel        string `mapstructure:"log_level"`
}

type HttpConfig struct {
	Addr           string        `mapstructure:"addr"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
	AllowIP        []string      `mapstructure:"allow_ip"`
}

type SwaggerConfig struct {
	Title    string   `mapstructure:"title"`
	Desc     string   `mapstructure:"desc"`
	Host     string   `mapstructure:"host"`
	BasePath string   `mapstructure:"base_path"`
	Schemes  []string `mapstructure:"schemes"`
	Version  string   `mapstructure:"version"`
}

func SetBaseConf(conf *Base) {
	Conf = conf
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

// GetStringConf 获取get配置信息
func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
	return confString
}

// GetStringMapConf 获取get配置信息
func GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetConf 获取get配置信息
func GetConf(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.Get(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetBoolConf 获取get配置信息
func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetFloat64Conf 获取get配置信息
func GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetIntConf 获取get配置信息
func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringMapStringConf 获取get配置信息
func GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringSliceConf 获取get配置信息
func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetTimeConf 获取get配置信息
func GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetDurationConf 获取时间阶段长度
func GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// IsSetConf 是否设置了key
func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.IsSet(strings.Join(keys[1:len(keys)], "."))
	return conf
}
