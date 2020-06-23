package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// Block 区块的基本结构与功能管理.
type Block struct {
	TimeStamp     int64          //时间戳
	Hash          []byte         //区块哈希值
	PrevBlockHash []byte         //前一个区块的哈希值
	Height        int64          //高度
	Txs           []*Transaction //交易数据
	Nonce         int64          //在运行pow时生成的哈希变化值，代表pow运行时动态修改的数据
}

//NewBlock 新建区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block Block

	block = Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Txs:           txs,
	}

	//生成哈希

	//block.SetHash()

	pow := NewProofOfWorf(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = int64(nonce)
	//执行工作量证明
	return &block
}

// SetHash 计算区块哈希值
// func (b *Block) SetHash() {

// 	timeStampBytes := IntToHex(b.TimeStamp)
// 	heightBytes := IntToHex(b.Height)

// 	blockBytes := bytes.Join([][]byte{
// 		heightBytes,
// 		timeStampBytes,
// 		b.PrevBlockHash,
// 		b.Txs,
// 	}, []byte{})
// 	//计算区块哈希
// 	hash := sha256.Sum256(blockBytes)
// 	b.Hash = hash[:]
// }

//CreateGenesisBlock 创建创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

//Serialize 区块结构数据化
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//gob
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if nil != err {
		log.Panicf("serialize the block in height %d to []byte failed %v\n", b.Height, err)
	}

	return buffer.Bytes()
}

//DeserialiazeBlock 区块数据反序列化
func DeserialiazeBlock(blockBytes []byte) *Block {
	var block Block
	//decoder
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panicf("deserialiaze the []byte to block failed %v\n", err)
	}
	return &block
}

//HashTransaction 把指定区块的所有交易结构都序列化
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte

	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}

	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
