package utils

import (
	"runtime"
	"path"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"reflect"
	"strconv"
)

func AppPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can not get current file info")
	}

	appPath := path.Dir(file) + "/../"

	return appPath
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串
	cipherStr := h.Sum(nil)
	result := fmt.Sprintf("%s", hex.EncodeToString(cipherStr)) // 输出加密结果

	return result
}

/**
 * 获取变量类型
 */
func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

/**
 * map to json
 */
func MapToJson(obj map[string]interface{}) string {
	jsonBytes, err := json.Marshal(obj)

	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return ""
	}

	return string(jsonBytes)
}

/**
 * json to map
 */
func JsonToMap(jsonStr string) map[string]interface{} {
	var mapResult map[string]interface{}

	if err := json.Unmarshal([]byte(jsonStr), &mapResult); err != nil {
		fmt.Println("json.Unmarshal failed:", err)
	}

	return mapResult
}

/**
 * int to string
 */
func IntToString(v int) string {
	return strconv.Itoa(v)
}

/**
 * string to int
 */
func StringToInt(str string) int {
	_int, _ := strconv.Atoi(str)
	return _int
}

/**
 * float64 to int
 */
func FloatToInt(v float64) int {
	return int(v)
}

/**
 * struct / interface to map
 */
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

/**
 * 字符串转数组
 */
func StringToRuneArr(str string, lens int) []rune {
	arr := []rune(str)

	if lens > len(arr) {
		panic("长度超出数组预期")
	}

	return arr[:lens]
}

/**
 * PHP-in_array()
 */
func InArray(value string, arr []string) bool {
	for _, val := range arr {
		if val == value {
			return true
		}
	}
	return false
}

