package blc

//BlockChain 区块链结构
type BlockChain struct {
	Blocks []*Block //区块的基本机构
}

//CreateBlockChainWithGenesis 初始化区块链
func CreateBlockChainWithGenesis() *BlockChain {
	block := CreateGenesisBlock([]byte("bc init"))
	return &BlockChain{[]*Block{block}}
}

//AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(height int64, prevBlockHash []byte, data []byte) {
	newBlock := NewBlock(height, prevBlockHash, data)
	bc.Blocks = append(bc.Blocks, newBlock)
}
