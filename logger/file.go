package logger

/**
将日志打印到文件中
 */
import (
	"fmt"
	"os"
	"strconv"
)

// 时间 级别 文件:行号
//格式化时间:2006-01-02 15:04:05.999
//定义文件日志所需要字段
type FileLogger struct {
	level       int
	logPath     string
	logName     string
	file        *os.File
	warnFile    *os.File
	LogDataChan chan *LogData
}

// 文件日志结构体构造方法
func NewFileLogger(config map[string]string) (log LogInterface, err error) {
	//判断配置中是否存在打印日志路径
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not fund log_path")
		return
	}
	//判断配置中是否存在日志文件名称
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not fund log_name")
		return
	}
	//判断配置中是否存在日志级别
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not fund log_level")
		return
	}
	//判断配置中是否存在日志级别
	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}

	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}
	//组建结构体文件日志对象
	log = &FileLogger{
		level:       getLogLevel(logLevel),
		logPath:     logPath,
		logName:     logName,
		LogDataChan: make(chan *LogData, chanSize),
	}
	//初始化文件对象
	log.Init()
	return
}

/**
初始化文件对象
 */
func (f *FileLogger) Init() {
	//组合文件名称和文件路径
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	//打开文件对象
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed,err:%v", filename, err))
	}
	//将文件对象存放文件日志结构体中
	f.file = file
	//写错误日志和fatal日志的文件
	errfilename := fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	errfile, err := os.OpenFile(errfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed,err:%v", filename, err))
	}
	f.warnFile = errfile
	go f.writeLogBackground()
}

/**
后台执行写入日志。
 */
func (f *FileLogger) writeLogBackground() {
	//这个for是阻塞循环，一直循环chan隧道
	for data := range f.LogDataChan {
		var file *os.File = f.file
		if data.WarnAndFatal {
			file = f.warnFile
		}
		fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", data.TimeStr, data.LevelStr, data.FileName, data.FuncName, data.LineNo,
			data.Message)
	}

}

/**
设置文件对象的日志级别
 */
func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

/**
打印dubag日志
 */
func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	// f.LogDataChan <- logData 将数据扔到队列里。
	// 通过select进行判断队列是否满了，如果满了就会走default分支；如果没满就将数据添加到chan中

	logData := epollWriteLog(LogLevelDebug, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	logData := epollWriteLog(LogLevelTrace, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := epollWriteLog(LogLevelInfo, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	logData := epollWriteLog(LogLevelWarn, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	logData := epollWriteLog(LogLevelError, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}

}
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	logData := epollWriteLog(LogLevelFatal, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

//关闭文件对象
func (f *FileLogger) Close() {
	f.warnFile.Close()
	f.file.Close()
}
