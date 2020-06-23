package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//交易管理文件

//Transaction 定义一个交易的基本结构
type Transaction struct {
	TxHash []byte      //交易哈希的标识
	Vins   []*TxInput  //输入列表
	Vouts  []*TxOutPut //输出列表
}

//NewCoinbaseTransaction 实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//输入
	txInput := &TxInput{[]byte{}, -1, "system reward"}
	//输出
	txOutput := &TxOutPut{10, address}
	txCoinbase := &Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutPut{txOutput},
	}

	txCoinbase.HashTransaction()
	return txCoinbase
}

//HashTransaction 生成交易哈希
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}

	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}
