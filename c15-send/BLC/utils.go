package blc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

//IntToHex 将int类型转换为[]byte。
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	//实现int64转byte[]
	return buffer.Bytes()
}

//JSONToSlice 标准的JSON格式转切片
func JSONToSlice(jsonString string) []string {
	var strSlice []string

	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Panicf("json to []string failed! %v\n", err)
	}

	return strSlice
}
