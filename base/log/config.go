package base

import (
	"errors"
	conf "iot_go/base/conf"
)

func BootstrapConf(fileName string) {
	/*初始化logger base*/
	defaultLoggerInit()
	log := &conf.Base{}
	conf.ParseConfigByFileName(fileName, log)
	conf.SetBaseConf(log)
	SetLogConf(log.LogConfig)
}

func SetLogConf(conf conf.LogConfig) error {
	if conf.CW.On {
		writer := NewConsoleWriter()
		writer.SetColor(conf.CW.Color)
		defaultLogger.registerLogWriter(writer)
	}
	if conf.FW.On {
		if len(conf.FW.LogPath) > 0 {
			writer := NewFileWriter()
			writer.SetFileName(conf.FW.LogPath)
			writer.SetPathPattern(conf.FW.RotateLogPath)
			writer.SetLogLevelFloor(TRACE)
			if len(conf.FW.WfLogPath) > 0 {
				writer.SetLogLevelCeil(INFO)
			} else {
				writer.SetLogLevelCeil(ERROR)
			}
			defaultLogger.registerLogWriter(writer)
		}

		if len(conf.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(conf.FW.WfLogPath)
			wfw.SetPathPattern(conf.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			defaultLogger.registerLogWriter(wfw)
		}
	}
	switch conf.Level {
	case "trace":
		defaultLogger.SetLevel(TRACE)
	case "debug":
		defaultLogger.SetLevel(DEBUG)
	case "info":
		defaultLogger.SetLevel(INFO)
	case "warning":
		defaultLogger.SetLevel(WARNING)
	case "error":
		defaultLogger.SetLevel(ERROR)
	case "fatal":
		defaultLogger.SetLevel(FATAL)
	default:
		return errors.New("invalid log_conf level")
	}
	return nil
}
