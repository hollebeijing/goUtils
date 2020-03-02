package iniconfig

import (
	"testing"
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

func TestIniConfig(t *testing.T) {
	//data, err := ioutil.ReadFile("config.ini")
	//if err != nil {
	//	t.Error("读取文件出错:", err)
	//}
	//var conf Config
	//err = UnMarshal(data, &conf)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//t.Logf("unmarshal success ,config :%#v", conf)
	//result, err := Marshal(conf)
	//if err != nil {
	//	t.Errorf("Marshal failed,err:%v", err)
	//
	//}
	//t.Logf("marshal success,config:%v", string(result))
	//
	//MarshalFile(conf, "./test.ini")

	var conf Config
	filename := "./test.ini"
	err := MarshalFile(filename, conf)
	if err != nil {
		t.Errorf("MarshalFile err:%v", err)
		return
	}
	var conf2 Config
	err = UnMarshalFile(filename, &conf2)
	if err != nil {
		t.Errorf("UnMarshalFile err:%v", err)
		return
	}
	t.Logf("UnMarshalFile %#v", conf2)
}
