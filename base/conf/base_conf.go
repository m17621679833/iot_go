package base

import (
	"github.com/spf13/viper"
	logFacede "iot_go/base/log"
)

var ViperConfMap map[string]*viper.Viper
var ConfBase *Conf

//var DBMapPool map[string]*sql.DB

// var GORMMapPool map[string]*gorm.DB
//var DBDefaultPool *sql.DB

// var GORMDefaultPool *gorm.DB
//var ConfRedis *RedisConf
//var ConfRedisMap *RedisMapConf

type Config struct {
	DebugMode    string `mapstructure:"debug_mode"`
	TimeLocation string `mapstructure:"time_location"`
}

type Conf struct {
	Log  logFacede.LogConfig `mapstructure:"log"`
	Base Config              `mapstructure:"base"`
}

type MysqlMapConf struct {
	List map[string]*MySQLConf `mapstructure:"list"`
}

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

type RedisMapConf struct {
	List map[string]*RedisConf `mapstructure:"list"`
}

type RedisConf struct {
	ProxyList    []string `mapstructure:"proxy_list"`
	Password     string   `mapstructure:"password"`
	Db           int      `mapstructure:"db"`
	ConnTimeout  int      `mapstructure:"conn_timeout"`
	ReadTimeout  int      `mapstructure:"read_timeout"`
	WriteTimeout int      `mapstructure:"write_timeout"`
}
