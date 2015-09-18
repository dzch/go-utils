package logger

import (
		"fmt"
		"os"
		"bufio"
		"runtime"
		"time"
		"sync"
		)

// log level
type LogLevel int
const (
		DEBUG = 16
		TRACE = 8
		NOTICE = 4
		WARNING = 2
		FATAL = 1
		OFF = 0
	  )
var level_prefix = map[LogLevel]string {
    DEBUG : "DEBUG:",
	TRACE : "TRACE:",
	NOTICE : "NOTICE:",
	WARNING : "WARNING:",
	FATAL : "FATAL:",
}

// logger struct
type _Logger struct {
	log_level LogLevel
	log_dir string
	log_file_prefix string
	nm_file_name string
	wf_file_name string
	nm_file *os.File
	wf_file *os.File
	nm_writer *bufio.Writer
	wf_writer *bufio.Writer
	mutex sync.Mutex
	pid int
}

var logger_obj *_Logger = nil
var time_format = "01-02 15:04:05:"

// check and re-create if needed
func (l *_Logger) _Check() (err error) {
	if l.nm_file == nil {
	    l.nm_file, err = os.OpenFile(l.nm_file_name, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), 0666)
	    if err != nil {
	    	l.nm_file.Close()
	    	return err
	    }
	    l.nm_writer = bufio.NewWriter(l.nm_file)
	} else {
		_, err := os.Stat(l.nm_file_name)
		if err != nil {
			l.nm_file.Close()
	        l.nm_file, err = os.OpenFile(l.nm_file_name, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), 0666)
	        if err != nil {
	        	l.nm_file.Close()
	        	return err
	        }
	        l.nm_writer = bufio.NewWriter(l.nm_file)
		}
	}
	if l.wf_file == nil {
	    l.wf_file, err = os.OpenFile(l.wf_file_name, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), 0666)
	    if err != nil {
	    	l.wf_file.Close()
	    	return err
	    }
	    l.wf_writer = bufio.NewWriter(l.wf_file)
	} else {
		_, err := os.Stat(l.wf_file_name)
		if err != nil {
			l.wf_file.Close()
	        l.wf_file, err = os.OpenFile(l.wf_file_name, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), 0666)
	        if err != nil {
	        	l.wf_file.Close()
	        	return err
	        }
	        l.wf_writer = bufio.NewWriter(l.wf_file)
		}
	}
	return nil
}

func (l *_Logger) _Write(level LogLevel, format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	err := l._Check()
	if err != nil {
		return
	}
	if level > WARNING {
	    l.nm_writer.WriteString(fmt.Sprintln(level_prefix[level], time.Now().Format(time_format), l.pid, fmt.Sprintf("%s:%d:", file, line), fmt.Sprintf(format, v...)))
	    l.nm_writer.Flush()
	} else {
	    l.wf_writer.WriteString(fmt.Sprintln(level_prefix[level], time.Now().Format(time_format), l.pid, fmt.Sprintf("%s:%d:", file, line), fmt.Sprintf(format, v...)))
	    l.wf_writer.Flush()
	}
}

func _NewLogger(dir string, file_prefix string, level LogLevel) (obj *_Logger, err error) {
	obj = &_Logger{
        log_level: level,
		log_dir : dir,
		log_file_prefix : file_prefix,
		nm_file_name : fmt.Sprintf("%s/%s.log", dir, file_prefix),
		wf_file_name : fmt.Sprintf("%s/%s.log.wf", dir, file_prefix),
		nm_file : nil,
		wf_file : nil,
		nm_writer : nil,
		wf_writer : nil,
	}
	return obj, nil
}

func Init(dir string, file_prefix string, level LogLevel) (err error) {
	if logger_obj != nil {
		return nil
	}
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		fmt.Println("logger.Init fail", err)
		return err
	}
    logger_obj, err = _NewLogger(dir, file_prefix, level)
	if err != nil {
		fmt.Println("logger.Init fail")
		return err
	}
	logger_obj.pid = os.Getpid()
	return nil
}

func Debug(format string, v ...interface{}) {
	if DEBUG > l.log_level {
		return
	}
	logger_obj._Write(DEBUG, format, v...)
}

func Trace(format string, v ...interface{}) {
	if TRACE > l.log_level {
		return
	}
	logger_obj._Write(TRACE, format, v...)
}

func Notice(format string, v ...interface{}) {
	if NOTICE > l.log_level {
		return
	}
	logger_obj._Write(NOTICE, format, v...)
}

func Warning(format string, v ...interface{}) {
	if WARNING > l.log_level {
		return
	}
	logger_obj._Write(WARNING, format, v...)
}

func Fatal(format string, v ...interface{}) {
	if FATAL > l.log_level {
		return
	}
	logger_obj._Write(FATAL, format, v...)
}

func Close() {
	if logger_obj == nil {
		return
	}
	if logger_obj.nm_file != nil {
	    logger_obj.nm_file.Close()
	}
	if logger_obj.wf_file != nil {
		logger_obj.wf_file.Close()
	}
}
