# goUtils

## 日志模块
### 使用方法
```$xslt
package main

import (
	"github.com/goUtils/logger"
	"fmt"
)

func initLogger(logPath, logName string, level string) (err error) {
	m := make(map[string]string, 8)
	m["log_path"] = logPath
	m["log_name"] = logName
	m["log_level"] = level

	err = logger.InitLogger("console", m)
	if err != nil {
		return
	}
	logger.Debug("init logger success")
	return
}

func Run()  {
	logger.Info("run info data")

}

func main()  {
	err := initLogger("./","test","debug")
	if err !=nil{
		fmt.Println("初始化日志系统错误:",err)
	}
	Run()
}
```