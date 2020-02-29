package logger

import "testing"

func TestFileLogger(t *testing.T)  {
	logger :=NewFileLogger(LogLevelDebug,"./","test")
	logger.Debug("user id[%d] is come from china",2222)
	logger.Warn("test warn log")
	logger.Fatal("test tatal log")
	logger.Close()
}


func TestConsoleLogger(t *testing.T)  {
	logger :=NewConsoleLogger(LogLevelDebug)
	logger.Debug("user id[%d] is come from china",2222)
	logger.Warn("test warn log")
	logger.Fatal("test tatal log")
	logger.Close()
}