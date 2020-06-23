package main

import (
	blc "bkc/c01-bc/BLC"
	"fmt"
)

func main() {
	block := blc.NewBlock(1, nil, []byte("the first block testing"))
	fmt.Printf("the first block : %v\n", block)
}
