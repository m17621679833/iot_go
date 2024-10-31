package base

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path"
	"time"
)

var pathVariableTable map[byte]func(*time.Time) int

type FileWriter struct {
	logLevelFloor int
	logLevelCeil  int
	fileName      string
	pathFormat    string
	file          *os.File
	fileBufWriter *bufio.Writer
	actions       []func(*time.Time) int
	variables     []interface{}
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (w *FileWriter) Init() error {
	return w.CreateLogFile()
}

func (w *FileWriter) CreateLogFile() error {
	if err := os.MkdirAll(path.Dir(w.fileName), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	if file, err := os.OpenFile(w.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		w.file = file
	}

	if w.fileBufWriter = bufio.NewWriterSize(w.file, 8192); w.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed")
	}
	return nil
}

func (w *FileWriter) Rotate() error {
	now := time.Now()
	v := 0
	rotate := false
	snapshot := make([]interface{}, len(w.variables))
	copy(snapshot, w.variables)
	for i, act := range w.actions {
		v = act(&now)
		if v != w.variables[i] {
			w.variables[i] = v
			rotate = true
		}
	}
	if rotate == false {
		return nil
	}
	return nil
}

func (w *FileWriter) SetFileName(name string) {
	w.fileName = name
}

func (w *FileWriter) SetLogLevelFloor(floor int) {
	w.logLevelFloor = floor
}

func (w *FileWriter) SetLogLevelCeil(ceil int) {
	w.logLevelCeil = ceil
}

func (w *FileWriter) SetPathPattern(pattern string) error {
	n := 0
	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}
	if n == 0 {
		w.pathFormat = pattern
		return nil
	}
	w.actions = make([]func(*time.Time) int, 0, n)
	w.variables = make([]interface{}, n)
	tmp := []byte(pattern)
	flag := false
	for _, c := range tmp {
		if flag == true {
			act, ok := pathVariableTable[c]
			if !ok {
				return errors.New("Invalid rotate pattern (" + pattern + ")")
			}
			w.actions = append(w.actions, act)
			flag = false
			continue
		}
		if c == '%' {
			flag = true
		}
	}

	for _, action := range w.actions {
		now := time.Now()
		w.variables = append(w.variables, action(&now))
	}
	w.pathFormat = convertPatternToFormat(tmp)

	return nil
}

func getYear(now *time.Time) int {
	return now.Year()
}

func getMonth(now *time.Time) int {
	return int(now.Month())
}

func getDay(now *time.Time) int {
	return now.Day()
}

func getHour(now *time.Time) int {
	return now.Hour()
}

func getMin(now *time.Time) int {
	return now.Minute()
}

func convertPatternToFormat(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	return string(pattern)
}

func init() {
	pathVariableTable = make(map[byte]func(*time.Time) int, 5)
	pathVariableTable['Y'] = getYear
	pathVariableTable['M'] = getMonth
	pathVariableTable['D'] = getDay
	pathVariableTable['H'] = getHour
	pathVariableTable['m'] = getMin
}

func (w *FileWriter) Write(r *LogRecord) error {
	if r.level < w.logLevelFloor || r.level > w.logLevelCeil {
		return nil
	}
	if w.fileBufWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

func (w *FileWriter) Flush() error {
	if w.fileBufWriter != nil {
		return w.fileBufWriter.Flush()
	}
	return nil
}
