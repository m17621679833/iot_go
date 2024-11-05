package base

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	config "iot_go/base/conf"
)

func GetGormPool(name string) (*gorm.DB, error) {
	if dbPool, ok := config.GORMMapPool[name]; ok {
		return dbPool, nil
	}
	return nil, errors.New("get pool error")
}

func InitSqlite3DB(path string) error {
	DbConfMap := &config.Sqlite3MapConf{}
	err := config.ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] empty db conf.")
	}
	config.GORMMapPool = map[string]*gorm.DB{}
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
		config.GORMMapPool[confName] = gormConn
	}
	if pool, err := GetSqlite3GormPool("default"); err == nil {
		config.GORMDefaultPool = pool
	}
	return nil
}

func GetSqlite3GormPool(name string) (*gorm.DB, error) {
	if pool, ok := config.GORMMapPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("get pool error")
}
