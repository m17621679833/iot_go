package base

import (
	"fmt"
	"sync"
)

var (
	LevelFlags = []string{"debug", "info", "warn", "error", "fatal", "panic"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const tunnelDefaultSize = 1024

type LogRecord struct {
	time  string
	code  string
	info  string
	level int
}

func (r *LogRecord) String() string {
	return fmt.Sprintf("[%s %s %s %s]\n", LevelFlags[r.level], r.time, r.code, r.info)
}

type LogWriter interface {
	Init() error
	Write(record *LogRecord) error
}

type LogRotate interface {
	Rotate() error
	SetPathPattern(pattern string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	logWriters    []LogWriter
	recordTunnel  chan *LogRecord
	level         int
	lastTimeStr   string
	lastTime      int64
	c             chan bool
	layout        string
	logRecordPool *sync.Pool
}

var (
	defaultLogger *Logger
	up            = false
)

func NewLogger() *Logger {
	if defaultLogger != nil && up == false {
		up = true
		return defaultLogger
	}
	logger := new(Logger)
	logger.logWriters = []LogWriter{}
	logger.recordTunnel = make(chan *LogRecord, tunnelDefaultSize)
	logger.c = make(chan bool, 2)
	logger.level = DEBUG
	logger.layout = "2006-01-02 15:04:05"
	logger.logRecordPool = &sync.Pool{New: func() interface{} {
		return &LogRecord{}
	}}
	go startLogWriterSilence(logger)
	return logger
}

func startLogWriterSilence(logger *Logger) {

}
