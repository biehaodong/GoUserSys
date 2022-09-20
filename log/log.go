package log

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	level   int
	logPath string
}

const (
	LevelFatal   int = iota // 严重错误信息.
	LevelError              // 错误信息.
	LevelWarning            // 警告信息.
	LevelInfo               // 普通信息.
	LevelDebug              // 调试信息.
)

var logger Logger

// ConfigLog: add log config.
func ConfigLog(logPath string, level int) error {
	// 打开日志文件.
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	logger.logPath = logPath
	logger.level = level
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logFile)

	return nil
}

// debug log.
func DebugLog(format string, v ...interface{}) {
	if logger.level >= LevelDebug {
		log.SetPrefix("debug ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// normal log.
func InfoLog(format string, v ...interface{}) {
	if logger.level >= LevelInfo {
		log.SetPrefix("info ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warning log.
func WarningLog(format string, v ...interface{}) {
	if logger.level >= LevelWarning {
		log.SetPrefix("warning ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error log .
func ErrorLog(format string, v ...interface{}) {
	if logger.level >= LevelError {
		log.SetPrefix("error ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal log.
func FatalLog(format string, v ...interface{}) {
	if logger.level >= LevelFatal {
		log.SetPrefix("fatal ")
		log.Output(2, fmt.Sprintf(format, v...))
	}
}
