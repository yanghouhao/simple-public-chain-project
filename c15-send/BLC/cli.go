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
	fmt.Printf("\tcreateblockchain -address address --create a blockchain\n")
	//添加区块
	fmt.Printf("\taddblock -data DATA --add block to blockchain\n")
	//打印完整的区块信息
	fmt.Printf("\tprintblockchain --print the information of blockchain\n")
	//通过命令行转账
	fmt.Printf("\tsend -from FROM -to TO -amount AMOUNT-- transfer AMOUNT from FROM to TO\n")
	fmt.Printf("\t\tthe descrition of transfer function\n")
	fmt.Printf("\t\t\t-from FROM -- the source address of this transaction\n")
	fmt.Printf("\t\t\t-to TO -- the destination address of this transaction\n")
	fmt.Printf("\t\t\t-amount AMOUNT -- the value of this transaction\n")
}

//createBlockchain 初始化区块链
func (cli *CLI) createBlockchain(address string) {
	CreateBlockChainWithGenesis(address)
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

//send 发起交易
func send() {
	if !dbExist() {
		fmt.Printf("there is no blockchain information\n")
		os.Exit(1)
	}
	//获取区块链对象
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock()
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
	//发起交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	//数据参数
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to yhh", "添加区块数据")
	//创建区块链时指定的矿工奖励
	flagCreateBlockchain := createBlockChainWithGenesisBlockCmd.String("address", "yhh", "system reward")
	//发起交易参数
	flagSendFromArg := sendCmd.String("from", "", "the source address of transaction")
	flagSendToArg := sendCmd.String("to", "", "the destination address of transaction")
	flagSendAmountArg := sendCmd.String("amount", "", "the value address of transaction")
	//判断命令
	switch os.Args[1] {
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse sendCmd failed! %v\n", err)
		}
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
		if *flagCreateBlockchain == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchain)
	}
	//./bc.exe send -from '[\"*\"]' -to '[\"aaa\"]' -amount '[\"10\"]'
	//发起转账
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Println("the source address shall not be nil")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg == "" {
			fmt.Println("the destination of transaction shall not be nil")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg == "" {
			fmt.Println("the value shall not be nil")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagSendAmountArg))
	}
}
