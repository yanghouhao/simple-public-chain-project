package blc

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

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

//dbExist 判断数据库文件是否存在
func dbExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		//数据库文件不存在
		return false
	}

	return true
}

//CreateBlockChainWithGenesis 初始化区块链
func CreateBlockChainWithGenesis(address string) *BlockChain {
	if dbExist() {
		//创世区块已存在
		fmt.Println("the genesis block has already been created")
		os.Exit(1)
	}
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
			//生成coinbase交易
			txCoinbase := NewCoinbaseTransaction(address)
			//生成创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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
func (bc *BlockChain) AddBlock(txs []*Transaction) {
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
			newBlock := NewBlock(latestBlock.Height+1, latestBlock.Hash, txs)
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

//PrintChain 遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain() {
	//读取数据库
	fmt.Println("the whole information of blockchain -----------------------------")
	var curBlock *Block

	bcit := bc.Iterator() //获取迭代器对象
	//循环读取
	//退出条件

	for {
		fmt.Println("-------------------------------------------")
		curBlock = bcit.Next()
		fmt.Printf("\tHash:%x\n", curBlock.Hash)
		fmt.Printf("\tPrevBlockHash:%x\n", curBlock.PrevBlockHash)
		fmt.Printf("\tTimeStamp:%v\n", curBlock.TimeStamp)
		fmt.Printf("\tHeight:%d\n", curBlock.Height)
		fmt.Printf("\tNonce:%d\n", curBlock.Nonce)
		fmt.Printf("\tTxs:%v\n", curBlock.Txs)

		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\t-hash : %x\n", tx.TxHash)
			fmt.Printf("\t\tinputs...\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-txHash : %x\n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout : %d\n", vin.Vout)
				fmt.Printf("\t\t\tvin-scriptSig : %s\n", vin.ScriptSig)
			}
			fmt.Printf("\t\toutputs...\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-value : %d\n", vout.Value)
				fmt.Printf("\t\t\tvout-scriptPubkey : %s\n", vout.ScriptPubkey)
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(curBlock.PrevBlockHash)
		//比较
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//遍历到创世区块
			break
		}
	}
}

//BlockChainObject 获取区块链对象
func BlockChainObject() *BlockChain {
	//获取DB
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("open the db [%s] faild! %v\n", dbName, err)
	}
	//获取Tip

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	if nil != err {
		log.Panicf("get the blockchain object failed %v\n", err)
	}

	return &BlockChain{DB: db, Tip: tip}
}

//MineNewBlock 通过接收交易生成区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	var block *Block
	var txs []*Transaction
	value, _ := strconv.Atoi(amount[0])
	//生成新的交易
	tx := NewSimpleTransaction(from[0], to[0], value)
	txs = append(txs, tx)
	//从数据库中获取最新区块
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserialiazeBlock(blockBytes)
		}
		return nil
	})
	//通过数据库中最新的区块去生成新区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	//持久化
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("update the new block to db failed %v \n", err)
			}

			err = b.Put([]byte("l"), block.Hash)
			if nil != err {
				log.Panicf("update the latest block hash to db failed %v \n", err)
			}
			bc.Tip = block.Hash
		}
		return nil
	})
}
