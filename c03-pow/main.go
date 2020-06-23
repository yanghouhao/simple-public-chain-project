package main

import (
	blc "bkc/c03-pow/BLC"
	"fmt"
)

func main() {
	// block := blc.NewBlock(1, nil, []byte("the first block testing"))
	// fmt.Printf("the first block : %v\n", block)
	bc := blc.CreateBlockChainWithGenesis()
	//fmt.Printf("blockchain:%v\n", bc.Blocks[0])
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("alice send 100 BTC to Bob"))

	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("Charlice send 100 BTC to Bob"))

	for _, block := range bc.Blocks {
		fmt.Printf("prevhash : %x blockhash : %x\n", block.PrevBlockHash, block.Hash)
	}
}
