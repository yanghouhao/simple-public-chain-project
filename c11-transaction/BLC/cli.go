package blc

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//对blockchain进行命令行管理

//CLI 对象
type CLI struct {
}

//PrintUsage 用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	//初始化区块链--
	fmt.Printf("\tcreateblockchain --create a blockchain\n")
	//添加区块
	fmt.Printf("\taddblock -data DATA --add block to blockchain\n")
	//打印完整的区块信息
	fmt.Printf("\tprintblockchain --print the information of blockchain\n")
}

//createBlockchain 初始化区块链
func (cli *CLI) createBlockchain(txs []*Transaction) {
	CreateBlockChainWithGenesis(txs)
}

//addBlock 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	//判断数据库是否存在
	if !dbExist() {
		fmt.Println("there is no blockchain....please call createblockchain command first")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.AddBlock(txs)
}

//printChain 打印
func (cli *CLI) printChain() {
	if !dbExist() {
		fmt.Println("there is no blockchain....please call createblockchain command first")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.PrintChain()
}

//IsValidArgs 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}
}

//Run 运行命令行
func (cli *CLI) Run() {
	IsValidArgs()
	//新建相关命令
	//添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	//输出区块链完整信息
	printChainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)
	//创建区块链
	createBlockChainWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	//数据参数
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to yhh", "添加区块数据")

	//判断命令
	switch os.Args[1] {
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse addBlockCmd failed! %v\n", err)
		}
	case "printblockchain":
		if err := printChainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse printChainCmd failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBlockChainWithGenesisBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBlockChainWithGenesisBlockCmd failed! %v\n", err)
		}
	default:
		//没有传递命令或者不在列表内
		PrintUsage()
		os.Exit(1)
	}

	//添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	//输出区块链信息
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	//创建区块链
	if createBlockChainWithGenesisBlockCmd.Parsed() {
		cli.createBlockchain([]*Transaction{})
	}
}
