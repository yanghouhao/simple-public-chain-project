package blc

import (
	"bytes"
	"crypto/sha256"
	"time"
)

// Block 区块的基本结构与功能管理.
type Block struct {
	TimeStamp     int64  //时间戳
	Hash          []byte //区块哈希值
	PrevBlockHash []byte //前一个区块的哈希值
	Height        int64  //高度
	Data          []byte //交易数据
}

//NewBlock 新建区块
func NewBlock(height int64, prevBlockHash []byte, data []byte) *Block {
	var block Block

	block = Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Data:          data,
	}

	//生成哈希

	block.SetHash()
	return &block
}

// SetHash 计算区块哈希值
func (b *Block) SetHash() {

	timeStampBytes := IntToHex(b.TimeStamp)
	heightBytes := IntToHex(b.Height)

	blockBytes := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		b.PrevBlockHash,
		b.Data,
	}, []byte{})
	//计算区块哈希
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

//CreateGenesisBlock 创建创世区块
func CreateGenesisBlock(data []byte) *Block {
	return NewBlock(1, nil, data)
}
