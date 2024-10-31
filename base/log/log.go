package base

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LevelFlags = []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
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
	return fmt.Sprintf("[%s] [%s] [%s] %s\n", r.time, LevelFlags[r.level], r.code, r.info)
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
	logWriters      []LogWriter
	recordTunnel    chan *LogRecord
	level           int
	lastTimeStr     string
	lastTime        int64
	terminalWatcher chan bool
	layout          string
	logRecordPool   *sync.Pool
}

func (l *Logger) productLogRecordToLogWriter(level int, format string, args ...interface{}) {
	var inf, code string
	if level < l.level {
		return
	}
	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}
	now := time.Now()
	if now.Unix() > l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := l.logRecordPool.Get().(*LogRecord)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level
	l.recordTunnel <- r
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
	logger.terminalWatcher = make(chan bool, 2)
	logger.level = DEBUG
	logger.layout = "2006-01-02 15:04:05"
	logger.logRecordPool = &sync.Pool{New: func() interface{} {
		return &LogRecord{}
	}}
	go boostrapLogWriterSilence(logger)
	up = true
	defaultLogger = logger
	return logger
}

func boostrapLogWriterSilence(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}
	var (
		r  *LogRecord
		ok bool
	)
	if r, ok = <-logger.recordTunnel; !ok {
		logger.terminalWatcher <- true
		return
	}
	log.Printf("boostrapLogWriterSilence %s\n", r.String())

	for _, writer := range logger.logWriters {
		if err := writer.Write(r); err != nil {
			log.Printf("write log writer failed: %v", err)
		}
	}
	flusherTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)
	for {
		select {
		case r, ok = <-logger.recordTunnel:
			if !ok {
				logger.terminalWatcher <- true
				return
			}
			log.Printf("ssss  %s\n", r.String())
			for _, writer := range logger.logWriters {
				if err := writer.Write(r); err != nil {
					log.Printf("write log writer failed: %v", err)
				}
			}
			logger.logRecordPool.Put(r)
		case <-flusherTimer.C:
			for _, writer := range logger.logWriters {
				if f, ok := writer.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Printf("flush log writer failed: %v", err)
					}
				}
			}
			flusherTimer.Reset(time.Millisecond * 1000)
		case <-rotateTimer.C:
			for _, writer := range logger.logWriters {
				if rotate, ok := writer.(LogRotate); ok {
					if err := rotate.Rotate(); err != nil {
						log.Printf("rotate log writer failed: %v", err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

func defaultLoggerInit() {
	if up == false {
		defaultLogger = NewLogger()
	}
}

func (l *Logger) SetLevel(level int) {
	defaultLogger.level = level
}
func (l *Logger) SetLayout(layout string) {
	defaultLogger.layout = layout
}

func (l *Logger) registerLogWriter(writer LogWriter) {
	if err := writer.Init(); err != nil {
		panic(err)
	}
	l.logWriters = append(l.logWriters, writer)
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(TRACE, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.productLogRecordToLogWriter(FATAL, format, args...)
}

func (l *Logger) Close() {
	close(l.recordTunnel)
	<-l.terminalWatcher
	for _, writer := range l.logWriters {
		if f, ok := writer.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Printf("flush log writer failed: %v", err)
			}
		}
	}
}
