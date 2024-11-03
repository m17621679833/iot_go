package base

import "errors"

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
			wfw.SetLogLevelFloor(WARNING)
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
		logger.SetLevel(WARNING)

	case "error":
		logger.SetLevel(ERROR)

	case "fatal":
		logger.SetLevel(FATAL)

	default:
		return errors.New("invalid log_conf level")
	}
	return nil
}
