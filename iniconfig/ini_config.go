package iniconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//将配置文件写入文件中
func MarshalFile(filename string, data interface{}) (err error) {
	result, err := Marshal(data)
	if err != nil {
		return
	}
	return ioutil.WriteFile(filename, result, 0755)

}

//将go中数据转换为ini格式数据
//这里使用data一个空接口，是因为我们目前不知道用户谁传什么数据进来的。
//序列化
func Marshal(data interface{}) (result []byte, err error) {
	//获取到数据的类型信息
	typeInfo := reflect.TypeOf(data)
	//获取到数据的值信息
	valueInfo := reflect.ValueOf(data)
	//判断传进来的数据是不是结构体，如果不是进行报错
	if typeInfo.Kind() != reflect.Struct {
		err = errors.New("please pass struct")
		return
	}
	var resData []string
	for i := 0; i < typeInfo.NumField(); i++ {
		// 每个区域字段
		sectionField := typeInfo.Field(i)
		// 获取每个区域字段的值
		sectionVal := valueInfo.Field(i)
		// 获取这个区域字段的类型信息
		sectionType := sectionField.Type
		// 判断当前类型是不是结构体，如果不是直接跳过
		if sectionType.Kind() != reflect.Struct {
			continue
		}
		// 类型信息获取tag中ini字段，如果没有就用当前字段名
		tagVal := sectionField.Tag.Get("ini")
		if len(tagVal) == 0 {
			tagVal = sectionField.Name
		}
		//组合数据将数据放到切片中
		section := fmt.Sprintf("\n[%s]\n", tagVal)
		resData = append(resData, section)

		// 获取当前字段的下面结构体的数据
		for j := 0; j < sectionType.NumField(); j++ {
			//找到key字段
			keyField := sectionType.Field(j)
			// 类型信息获取tag中ini字段，如果没有就用当前字段名
			fieldTagVal := keyField.Tag.Get("ini")
			if len(fieldTagVal) == 0 {
				fieldTagVal = keyField.Name
			}
			//找到val数据
			valField := sectionVal.Field(j)
			//组合数据将数据放到切片中
			item := fmt.Sprintf("%s=%v\n", fieldTagVal, valField.Interface())
			resData = append(resData, item)

		}
	}
	for _, v := range resData {
		result = append(result, []byte(v)...)
	}
	return
}

//将文件中配置读取到go的程序中
func UnMarshalFile(filename string, result interface{}) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return UnMarshal(data, result)
}

//反序列化,将ini格式数据转换为go的数据
/*
思路:
	1、先将传进来的配置文件内容进行解析
	2、将解析后的数据一一对应到结构体上

 */
func UnMarshal(data []byte, result interface{}) (err error) {
	lineArr := strings.Split(string(data), "\n")

	typeInfo := reflect.TypeOf(result)
	if typeInfo.Kind() != reflect.Ptr {
		err = errors.New("please pass address")
		return
	}

	//这里将指针类型变成值类型
	typeStruct := typeInfo.Elem()
	if typeStruct.Kind() != reflect.Struct {
		err = errors.New("please pass struct")
	}

	var lastFiledName string

	for index, line := range lineArr {
		//这里将每一个行前后空格去掉
		line = strings.TrimSpace(line)
		//这里如果是空行直接跳过
		if len(line) == 0 {
			continue
		}
		//如果是注释，直接跳过
		if line[0] == ';' || line[0] == '#' {
			continue
		}
		if line[0] == '[' {
			lastFiledName, err = parseSection(line, typeStruct)
			if err != nil {
				err = fmt.Errorf("%v lineNo:%d", err, index+1)
				return
			}
			continue
		}

		if len(lastFiledName) == 0 {
			err = fmt.Errorf("not Section lineNo:%d", index+1)
			return
		}
		err = parseItem(lastFiledName, line, result)
		if err != nil {
			err = fmt.Errorf("%v lineNo:%d", err, index+1)
			return
		}

	}
	return
}

func parseItem(filedName, line string, result interface{}) (err error) {
	index := strings.Index(line, "=")
	if index == -1 {
		err = fmt.Errorf("syntax error,line:%s", line)
		return
	}

	key := strings.TrimSpace(line[0:index])
	val := strings.TrimSpace(line[index+1:])
	if len(key) == 0 {
		err = fmt.Errorf("syntax error,line:%s", line)
		return
	}
	resultValue := reflect.ValueOf(result)
	sectionValue := resultValue.Elem().FieldByName(filedName) //这个也是一个结构体，
	sectionType := sectionValue.Type()
	if sectionType.Kind() != reflect.Struct {
		err = fmt.Errorf("field:%s must be struct", filedName)
		return
	}
	keyFiledName := ""
	for i := 0; i < sectionType.NumField(); i++ {
		filed := sectionType.Field(i)
		tagVal := filed.Tag.Get("ini")
		if tagVal == key {
			keyFiledName = filed.Name
			break
		}
	}
	if len(keyFiledName) == 0 {
		return
	}

	filedValue := sectionValue.FieldByName(keyFiledName)
	if filedValue == reflect.ValueOf(nil) {
		return
	}
	fieldKind := filedValue.Type().Kind()
	switch fieldKind {
	case reflect.String:
		filedValue.SetString(val)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		intVal, errRet := strconv.ParseInt(val, 10, 64)
		if errRet != nil {
			err = errRet
			return
		}
		filedValue.SetInt(intVal)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		intVal, errRet := strconv.ParseUint(val, 10, 64)
		if errRet != nil {
			err = errRet
			return
		}
		filedValue.SetUint(intVal)
	case reflect.Float32, reflect.Float64:
		intVal, errRet := strconv.ParseFloat(val, 64)
		if errRet != nil {
			err = errRet
			return
		}
		filedValue.SetFloat(intVal)
	case reflect.Bool:
		intVal, errRet := strconv.ParseBool(val)
		if errRet != nil {
			err = errRet
			return
		}
		filedValue.SetBool(intVal)
	default:
		err = fmt.Errorf("unsupport type:%v", fieldKind)
	}
	return

}

//获取【server】外层配置名
func parseSection(line string, typeInfo reflect.Type, ) (fileName string, err error) {
	//判断合法性
	if line[0] == '[' && len(line) <= 2 {
		err = fmt.Errorf("syntax error,invalid section:%s", line)
		return
	}
	if line[0] == '[' && line[len(line)-1] != ']' {
		err = fmt.Errorf("syntax error,invalid section:%s", line)
		return
	}
	if line[0] == '[' && line[len(line)-1] == ']' {
		sectionName := strings.TrimSpace(line[1 : len(line)-1])
		if len(sectionName) == 0 {
			err = fmt.Errorf("syntax error,invalid section:%s", line)
			return
		}
		for i := 0; i < typeInfo.NumField(); i++ {
			filed := typeInfo.Field(i)
			tagValue := filed.Tag.Get("ini")
			if tagValue == sectionName {
				fileName = filed.Name
				break
			}
		}
	}
	return
}
