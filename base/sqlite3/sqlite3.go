package base

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	fileParser "iot_go/base/conf"
	"time"
)

var GORMMapPool map[string]*gorm.DB
var GORMDefaultPool *gorm.DB

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

func GetGormPool(name string) (*gorm.DB, error) {
	if dbPool, ok := GORMMapPool[name]; ok {
		return dbPool, nil
	}
	return nil, errors.New("get pool error")
}

func InitSqlite3DB(path string) error {
	DbConfMap := &Sqlite3MapConf{}
	err := fileParser.ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(fileParser.TimeFormat), " empty sqlite3 conf.")
	}
	GORMMapPool = map[string]*gorm.DB{}
	for confName, conf := range DbConfMap.List {
		gormConn, err := gorm.Open(sqlite.Open(conf.DataSourceName), &gorm.Config{})
		if err != nil {
			return err
		}
		gormConn.Exec("PRAGMA journal_mode=WAL;")
		gormConn.Exec("PRAGMA synchronous = OFF;")
		if gormConn.Error != nil {
			panic("failed to set WAL mode")
		}
		gormDb, err := gormConn.DB()
		if err != nil {
			return err
		}
		gormDb.SetMaxIdleConns(conf.MaxIdleConn)
		gormDb.SetMaxOpenConns(conf.MaxOpenConn)
		err = gormDb.Ping()
		if err != nil {
			return err
		}
		GORMMapPool[confName] = gormConn
	}
	if pool, err := GetSqlite3GormPool("default"); err == nil {
		GORMDefaultPool = pool
	}
	return nil
}

func GetSqlite3GormPool(name string) (*gorm.DB, error) {
	if pool, ok := GORMMapPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("get pool error")
}
