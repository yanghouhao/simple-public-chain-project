package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

//区块链迭代器管理文件

//BlockChainIterator 实现迭代器基本结构
type BlockChainIterator struct {
	DB      *bolt.DB //迭代目标
	CurHash []byte   //当前目标的哈希
}

//Iterator 创建迭代器对象
func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{blc.DB, blc.Tip}
}

//Next 实现迭代函数next，获取每一个区块
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block

	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			curBlockBytes := b.Get(bcit.CurHash)
			block = DeserialiazeBlock(curBlockBytes)
			//更新区块哈希值
			bcit.CurHash = block.PrevBlockHash
		}
		return nil
	})

	if nil != err {
		log.Panicf("iterator the db failed! %v\n", err)
	}
	return block
}
