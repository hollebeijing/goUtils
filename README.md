# goUtils

## 日志模块
### 使用方法
```$xslt
package main

import (
	"fmt"
	"github.com/goUtils/logger"
)

func initLogger(logPath, logName string, level string) (err error) {
	m := make(map[string]string, 8)
	m["log_path"] = logPath
	m["log_name"] = logName
	m["log_level"] = level
	m["log_split_type"] ="size" //配置大小分隔,默认为每小时切分一次
	m["log_split_size"] ="104857600" // 配置多大进行切分，默认104857600(10M)


	err = logger.InitLogger("file", m)
	if err != nil {
		return
	}
	logger.Debug("init logger success")
	return
}

func Run()  {
	for {
		logger.Info("run info data")
		//time.Sleep(time.Second)
	}


}

func main()  {
	err := initLogger("./","test","debug")
	if err !=nil{
		fmt.Println("初始化日志系统错误:",err)
	}
	Run()
}


```

## 配置文件模块
### 使用说明
```$xslt
package main

import (
	"fmt"
	"github.com/goUtils/iniconfig"
)

type Config struct {
	ServerConf ServerConfig `ini:"server"`
	MysqlConf  MysqlConfig  `ini:"mysql"`
}

type ServerConfig struct {
	Ip   string `ini:"ip"`
	Port int    `ini:"port"`
}

type MysqlConfig struct {
	Auth     bool   `ini:"auth"`
	Username string `ini:"username"`
	Passwd   string `ini:"passwd"`
	DB       string `ini:"db"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

func main() {
	filename := "./config.ini"
	var conf Config
	err := iniconfig.UnMarshalFile(filename, &conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("config:%#v", conf)
}

```