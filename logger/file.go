package logger

/**
将日志打印到文件中
 */
import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// 时间 级别 文件:行号
//格式化时间:2006-01-02 15:04:05.999
//定义文件日志所需要字段
type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	LogDataChan   chan *LogData
	logSplitType  int
	logSplitSize  int64
	lastSplitHour int
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
	var logSplitType int = LogSplitTypeHour
	var logSplitSize int64
	//判断配置中是否存在日志级别
	logSplitStr, ok := config["log_split_type"]
	if !ok {
		logSplitType = LogSplitTypeHour
	} else {
		if logSplitStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857600" // 104857600 = 100M
			}
			// 第一个参数:要转换的字符数字，第二个传参数:要转换的进制:10进制，第三个参数:位数:64位
			logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeHour
		}
	}

	//判断配置中是否存在异步隧道大小，如果没有默认为长度为50000
	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}
	//这里传过来的都是字符串，因此这里要用strconv.Atoi 将字符的数字转为int格式
	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}

	//组建结构体文件日志对象
	log = &FileLogger{
		level:         getLogLevel(logLevel),
		logPath:       logPath,
		logName:       logName,
		LogDataChan:   make(chan *LogData, chanSize),
		logSplitSize:  logSplitSize,
		logSplitType:  logSplitType,
		lastSplitHour: time.Now().Hour(),
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
		f.checkSplitFile(data.WarnAndFatal)
		fmt.Fprintf(file, "%s %s (%s:%s:%d) %s\n", data.TimeStr, data.LevelStr, data.FileName, data.FuncName, data.LineNo,
			data.Message)
	}

}

/**
1、获取当前时间进行对比

 */
func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	// 如果最后存储的小时和当前是一个小时,那个不需要切分，如果不是则下面进行切分
	if (hour == f.lastSplitHour) {
		return
	}
	// 获取到文件对象
	file := f.file
	//定义一下源文件名，和备份文件名称
	var backupFileNmae, fileName string
	// 先判断是要切分哪个文件
	if warnFile {
		//切分错误日志文件
		backupFileNmae = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d", f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
		file = f.warnFile
	} else {
		//切分正常日志文件
		backupFileNmae = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d", f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	// 关闭文件
	file.Close()
	//将文件修改名称
	os.Rename(fileName, backupFileNmae)

	//创建新文件

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}

}

func (f *FileLogger) splitFileSize(warnFile bool) {
	now := time.Now()
	// 获取到文件对象
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	// 获取到文件基本信息
	statInfo, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	//定义一下源文件名，和备份文件名称
	var backupFileNmae, fileName string
	// 先判断是要切分哪个文件
	if warnFile {
		//切分错误日志文件
		backupFileNmae = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d%02d%02d", f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
		file = f.warnFile
	} else {
		//切分正常日志文件
		backupFileNmae = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d%02d%02d", f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	// 关闭文件
	file.Close()
	//将文件修改名称
	os.Rename(fileName, backupFileNmae)

	//创建新文件
	newfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = newfile
	} else {
		f.file = newfile
	}

}

func (f *FileLogger) checkSplitFile(warnFile bool) {
	if (f.logSplitType == LogSplitTypeHour) {
		f.splitFileHour(warnFile)
		return
	}
	f.splitFileSize(warnFile)

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
