package util

import (
	"bytes"
	"encoding/gob"
)

//对象的深度拷贝
//结构中的变量必须首字母大写才能被拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
