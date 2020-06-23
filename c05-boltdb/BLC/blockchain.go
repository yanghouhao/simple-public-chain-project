package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

//区块链管理文件

//数据库名称
const dbName = "block.db"

//表名称
const blockTableName = "blocks"

//BlockChain 区块链结构
type BlockChain struct {
	//Blocks []*Block //区块的基本机构
	DB  *bolt.DB //	数据库对象
	Tip []byte   //保存最新区块的哈希值
}

//CreateBlockChainWithGenesis 初始化区块链
func CreateBlockChainWithGenesis() *BlockChain {
	//保存最新区块的哈希值
	var blockHash []byte
	//创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("create db [%s] failed %v\n", dbName, err)
	}
	//创建桶,把创世区块存入数据库
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			//没找到
			buc, err := tx.CreateBucket([]byte(blockTableName))
			if nil != err {
				log.Panicf("create bucket [%s] failed %v \n", blockTableName, err)
			}
			//生成创世区块
			genesisBlock := CreateGenesisBlock([]byte("bc init"))
			//存储 用hash作为key,序列化作为value
			err = buc.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert genisis block failed %v\n", err)
			}
			blockHash := genesisBlock.Hash
			//存储对应区块的哈希
			err = buc.Put([]byte("l"), blockHash)
			if nil != err {
				log.Panicf("save the latest hash of genisis block failed %v\n", err)
			}
		}
		return nil
	})
	return &BlockChain{DB: db, Tip: blockHash}
}

//AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(data []byte) {
	//更新区块（insert）
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			//获取最新区块的哈希值
			//fmt.Printf("latest hash %v\n", b.Get([]byte("l")))
			blockBytes := b.Get(b.Get([]byte("l")))
			//反序列化
			latestBlock := DeserialiazeBlock(blockBytes)
			//新建区块
			newBlock := NewBlock(latestBlock.Height+1, latestBlock.Hash, data)
			//存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if nil != err {
				log.Panicf("insert new block to db failed %v\n", err)
			}
			//更新最新区块的哈希值
			err = b.Put([]byte("l"), newBlock.Hash)
			if nil != err {
				log.Panicf("update the latest block in Height : %d to db failed %v\n", newBlock.Height, err)
			}

			//更新区块链中的最新哈希
			bc.Tip = newBlock.Hash
		}
		return nil
	})

	if nil != err {
		log.Panicf("insert block to db failed %v", err)
	}
}
