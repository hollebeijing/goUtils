package logger

//定义日志接口规范
type LogInterface interface {
	//定义日志级别
	SetLevel(level int)
	//日志级别
	Debug(format string,args ...interface{})
	Trace(format string,args ...interface{})
	Info(format string,args ...interface{})
	Warn(format string,args ...interface{})
	Error(format string,args ...interface{})
	Fatal(format string,args ...interface{})
	//关闭文件
	Close()
	//初始化文件
	Init()

}