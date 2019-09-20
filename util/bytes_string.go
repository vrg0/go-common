package util

import "unsafe"

//注意，此方法需要保证原data存在，一般转换后的字符串应该作为临时变量使用
func BytesString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
