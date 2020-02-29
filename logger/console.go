package logger
/**
终端打印日志
将日志打印到终端中
 */
import (
	"fmt"
	"os"
)

//定义终端日志结构体
type ConsoleLogger struct {
	level int
}
// 终端日志结构体构造方法
func NewConsoleLogger(config map[string]string) (log LogInterface, err error) {
	//判断日志是否存在日志级别
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not fund log_level")
		return
	}
	log = &ConsoleLogger{
		//用户传过来的日志字符串小写，转换为数字常量
		level: getLogLevel(logLevel),
	}
	return
}


// 设置日志级别
func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	c.level = level
}
/**
打印dubag日志
 */
func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}
	//真实打印日志公共方法
	//os.Stdout:go中终端也是一个文件对象，直接传入即可。
	writeLog(os.Stdout, LogLevelDebug, format, args...)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	writeLog(os.Stdout, LogLevelTrace, format, args...)

}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	writeLog(os.Stdout, LogLevelInfo, format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	writeLog(os.Stdout, LogLevelWarn, format, args...)

}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	writeLog(os.Stdout, LogLevelError, format, args...)

}
func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	writeLog(os.Stdout, LogLevelFatal, format, args...)
}

func (c *ConsoleLogger) Close() {
}
func (c *ConsoleLogger) Init() {
}

