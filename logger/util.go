package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

/**
使用runtime函数获取 go文件、方法名、行号
 */
func GetLineInfo() (fileName, funcName string, lineNo int) {
	//runtime主要获取go文件 方法名，行号数据
	//skip是堆栈层级，0就是当前位置的方法
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		//获取文件名称绝对路径
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}

//真实打印日志公共方法
func writeLog(file *os.File, level int, format string, args ...interface{}) {
	// 获取当前时间
	now := time.Now()
	//格式化字符串时间
	nowStr := now.Format("2006-01-02 15:04:05.999")
	//获取日志级别大写字符串打印到日志中
	levelStr := getLevelText(level)
	//获取文件名，方法名 行号
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	//组合格式化用户要打印的数据
	msg := fmt.Sprintf(format, args...)
	//组合整体日志中打印的内容,并打印
	fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", nowStr, levelStr, fileName, funcName, lineNo,
		msg)
}

/**
1、当业务调用打印日志的方法时，我们把日志相关的数据写入到chan(队列)
2、然后我们有一个后台的线程不断的从chan里面获取这些日志，最终写入到日志中。
 */

type LogData struct {
	Message      string
	TimeStr      string
	LevelStr     string
	FileName     string
	FuncName     string
	LineNo       int
	WarnAndFatal bool
}

//真实异步打印日志公共方法
func epollWriteLog(level int, format string, args ...interface{}) *LogData {
	// 获取当前时间
	now := time.Now()
	//格式化字符串时间
	nowStr := now.Format("2006-01-02 15:04:05.999")
	//获取日志级别大写字符串打印到日志中
	levelStr := getLevelText(level)
	//获取文件名，方法名 行号
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	//组合格式化用户要打印的数据
	msg := fmt.Sprintf(format, args...)
	//组合整体日志中打印的内容,并打印
	//fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", nowStr, levelStr, fileName, funcName, lineNo,
	//	msg)
	WarnAndFatal := false
	if level == LogLevelError || level == LogLevelWarn || level == LogLevelFatal {
		WarnAndFatal = true
	}
	logData := &LogData{
		msg,
		nowStr,
		levelStr,
		fileName,
		funcName,
		lineNo,
		WarnAndFatal,
	}

	return logData
}
