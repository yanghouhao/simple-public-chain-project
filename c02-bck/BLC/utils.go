package blc

import (
	"bytes"
	"encoding/binary"
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
