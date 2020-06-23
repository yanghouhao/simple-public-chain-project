package blc

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

//共识管理文件

//实现POW实例以及相关功能

//目标难度
const targetBit = 16

//ProofOfWork 工作量证明的结构
type ProofOfWork struct {
	Block  *Block   //需要共识的区块
	target *big.Int //目标难度哈希
}

//NewProofOfWorf 创建一个POW对象
func NewProofOfWorf(block *Block) *ProofOfWork {
	//str1 hash
	//strTarget
	target := big.NewInt(1)
	//需要前两位为零才能满足要求
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

//Run 返回哈希值，返回碰撞次数执行pow 比较哈希
func (proofOfWork *ProofOfWork) Run() ([]byte, int) {
	//碰撞次数
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte
	for {
		//生成准备数据
		dataBytes := proofOfWork.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//检测生成的哈希值是否符合条件
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			//找到了符合条件的哈希，
			break
		}
		nonce++
	}

	//fmt.Printf("打印碰撞次数  ： %d\n", nonce)

	return hash[:], nonce
}

//prepareData 执行准备工作
func (proofOfWork *ProofOfWork) prepareData(nonce int64) []byte {
	var data []byte
	//拼接区块属性
	timeStampBytes := IntToHex(proofOfWork.Block.TimeStamp)
	heightBytes := IntToHex(proofOfWork.Block.Height)
	data = bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		proofOfWork.Block.PrevBlockHash,
		proofOfWork.Block.Data,
		IntToHex(targetBit),
		IntToHex(nonce),
	}, []byte{})

	return data
}
