package base

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logger := NewLogger()
	conf := LogConfig{
		Level: "warn",
		FW: FileWriterConf{
			On:              true,
			LogPath:         "/log_test.log_conf",
			RotateLogPath:   "./log_test.log_conf",
			WfLogPath:       "./log_test.wf.log_conf",
			RotateWfLogPath: "./log_test.wf.log_conf",
		},
		CW: ConsoleWriterConf{
			On:    true,
			Color: true,
		},
	}
	SetLogConf(conf)
	logger.Warn("test message")
	logger.Close()
	time.Sleep(time.Second)
}
