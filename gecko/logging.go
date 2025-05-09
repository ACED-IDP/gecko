package gecko

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/uc-cdis/arborist/arborist"
)

const (
	LogLevelDebug   arborist.LogLevel = "DEBUG"
	LogLevelInfo    arborist.LogLevel = "INFO"
	LogLevelWarning arborist.LogLevel = "WARNING"
	LogLevelError   arborist.LogLevel = "ERROR"
)

type LogCache struct {
	logs []Log
}

type Log struct {
	lvl arborist.LogLevel
	msg string
}

func (cache *LogCache) write(logger arborist.Logger) {
	for _, log := range cache.logs {
		logger.Print(log.msg)
	}
}

func (handler *LogHandler) Print(format string, a ...interface{}) {
	handler.logger.Print(sprintf(format, a...))
}

func (handler *LogHandler) Debug(format string, a ...interface{}) {
	handler.logger.Print(logMsg(LogLevelDebug, format, a...))
}

func (handler *LogHandler) Info(format string, a ...interface{}) {
	handler.logger.Print(logMsg(LogLevelInfo, format, a...))
}

func (handler *LogHandler) Warning(format string, a ...interface{}) {
	handler.logger.Print(logMsg(LogLevelWarning, format, a...))
}

func (handler *LogHandler) Error(format string, a ...interface{}) {
	handler.logger.Print(logMsg(LogLevelError, format, a...))
}

func (cache *LogCache) Debug(format string, a ...interface{}) {
	log := Log{
		lvl: LogLevelDebug,
		msg: logMsg(LogLevelDebug, format, a...),
	}
	cache.logs = append(cache.logs, log)
}

func (cache *LogCache) Info(format string, a ...interface{}) {
	log := Log{
		lvl: LogLevelInfo,
		msg: logMsg(LogLevelInfo, format, a...),
	}
	cache.logs = append(cache.logs, log)
}

func (cache *LogCache) Warning(format string, a ...interface{}) {
	log := Log{
		lvl: LogLevelWarning,
		msg: logMsg(LogLevelWarning, format, a...),
	}
	cache.logs = append(cache.logs, log)
}

func (cache *LogCache) Error(format string, a ...interface{}) {
	log := Log{
		lvl: LogLevelError,
		msg: logMsg(LogLevelError, format, a...),
	}
	cache.logs = append(cache.logs, log)
}

func logMsg(lvl arborist.LogLevel, format string, a ...interface{}) string {
	msg := sprintf(format, a...)
	msg = fmt.Sprintf("%s: %s", lvl, msg)
	// get the call from 2 stack frames above this
	// (one call up is the LogCache method, so go one more above that)
	_, fn, line, ok := runtime.Caller(2)
	if ok {
		// shorten the filepath to only the basename
		split := strings.Split(fn, string(os.PathSeparator))
		fn = split[len(split)-1]
		msg = fmt.Sprintf("%s:%d: %s", fn, line, msg)
	}
	return msg
}

func sprintf(format string, a ...interface{}) string {
	var msg string
	if len(a) == 0 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	return msg
}
