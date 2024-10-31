package base

import "errors"

type FileWriterConf struct {
	On              bool
	LogPath         string
	RotateLogPath   string
	WfLogPath       string
	RotateWfLogPath string
}

type ConsoleWriterConf struct {
	On    bool
	Color bool
}

type LogConfig struct {
	Level string
	FW    FileWriterConf
	CW    ConsoleWriterConf
}

func SetLogConf(conf LogConfig) {
	defaultLoggerInit()
	InstanceLogConf(conf, defaultLogger)
}

func InstanceLogConf(conf LogConfig, logger *Logger) error {
	if conf.CW.On {
		writer := NewConsoleWriter()
		writer.SetColor(conf.CW.Color)
		logger.registerLogWriter(writer)
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
			logger.registerLogWriter(writer)
		}

		if len(conf.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(conf.FW.WfLogPath)
			wfw.SetPathPattern(conf.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARN)
			wfw.SetLogLevelCeil(ERROR)
			logger.registerLogWriter(wfw)
		}
	}
	switch conf.Level {
	case "trace":
		logger.SetLevel(TRACE)

	case "debug":
		logger.SetLevel(DEBUG)

	case "info":
		logger.SetLevel(INFO)

	case "warn":
		logger.SetLevel(WARN)

	case "error":
		logger.SetLevel(ERROR)

	case "fatal":
		logger.SetLevel(FATAL)

	default:
		return errors.New("invalid log_conf level")
	}
	return nil
}
