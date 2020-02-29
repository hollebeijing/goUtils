package logger
/**
封装面对用户接口，让用户更方便的使用。
 */
var log LogInterface

/*
name:
	file:"初始化一个文件日志实例"
	console:"初始化console日志实例"
config:
	log_path:文件路径,打印到终端可以没有
	log_name:文件名称,打印到终端可以没有
	log_level:日志级别
 */
func InitLogger(name string, config map[string]string) (err error) {
	switch name {
	case "file":
		log, err = NewFileLogger(config)
	case "console":
		log, err = NewConsoleLogger(config)
	default:
		log, err = NewFileLogger(config)
	}
	return
}

//封装对外访问数据
func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}
func Trace(format string, args ...interface{}) {
	log.Trace(format, args...)
}
func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}
func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}
func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}
func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}
