package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	config "iot_go/base/conf"
	base "iot_go/base/log"
	"iot_go/base/util"
	"time"
)

func transferLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

func GetDefaultDB(g *gin.Context) (*gorm.DB, error) {
	if dbPool, ok := config.GORMMapPool["default"]; ok {
		traceContext := util.GetGinTraceContext(g)
		dbPool = dbPool.Set("trace_context", traceContext)
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
		level := transferLevel(conf.LogLevel)
		gormConn, err := gorm.Open(sqlite.Open(conf.DataSourceName), &gorm.Config{Logger: &GormLogger{
			LogLevel:   level,
			NiceLogger: nil,
		}})
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

type GormLogger struct {
	LogLevel   logger.LogLevel
	NiceLogger *base.Logger
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录 Info 级别的日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.NiceLogger.Info(msg, data...)
		fmt.Printf("[INFO] %s: %v\n", msg, data)
	}
}

// Warn 记录 Warn 级别的日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		fmt.Printf("[WARN] %s: %v\n", msg, data)
	}
}

// Error 记录 Error 级别的日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		fmt.Printf("[ERROR] %s: %v\n", msg, data)
	}
}

// Trace 记录 SQL 执行的详细信息
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		fmt.Printf("[ERROR] %s [%s] %v\n", sql, elapsed, err)
	case elapsed > time.Second && l.LogLevel >= logger.Warn:
		fmt.Printf("[WARN] %s [%s] (slow query)\n", sql, elapsed)
	case l.LogLevel >= logger.Info:
		fmt.Printf("[INFO] %s [%s] %d rows affected\n", sql, elapsed, rows)
	}
}
