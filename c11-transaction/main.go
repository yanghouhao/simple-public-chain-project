package main

/*
 * @Author: your name
 * @Date: 2020-06-23 14:19:12
 * @LastEditTime: 2020-06-23 14:20:42
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \bkc\c11-transaction\main.go
 */

import (
	blc "bkc/c11-transaction/BLC"
)

func main() {

	// bc := blc.CreateBlockChainWithGenesis()
	// bc.AddBlock([]byte("a send 100 btc to b"))
	// bc.AddBlock([]byte("b send 100 btc to c"))
	// bc.AddBlock([]byte("a send 100 btc to c"))

	// bc.PrintChain()
	cli := blc.CLI{}
	cli.Run()
}
