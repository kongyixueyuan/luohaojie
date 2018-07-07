package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}


func printUsage()  {

	fmt.Println("使用说明:")


	fmt.Println("\tcreateblockchain -address -- 创建创世区块.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")
	fmt.Println("\taddresslists -- 输出所有钱包地址.")
	fmt.Println("\tcreatewallet -- 创建钱包.")

}

//过滤参数数量
func isValidArgs()  {
	//获取当前输入参数个数
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}



func (cli *CLI) Run()  {

	isValidArgs()


	//自定义cli命令
	addresslistsCmd := flag.NewFlagSet("addresslists",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)


	flagFrom := sendBlockCmd.String("from","","源地址")
	flagTo := sendBlockCmd.String("to","","目的地址")
	flagAmount := sendBlockCmd.String("amount","","转账金额")


	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address","","查询账号余额")



	switch os.Args[1] {
	//第二个参数为相应命令，取第三个参数开始作为参数并解析
		case "send":
			err := sendBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "addresslists":
			err := addresslistsCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "getbalance":
			err := getbalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createwallet":
			err := createWalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			printUsage()
			os.Exit(1)
	}
	//对addBlockCmd命令的解析
	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}



		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		//输入地址有效性判断
		for index,fromAdress := range from {
			if IsValidForAddress([]byte(fromAdress)) == false || IsValidForAddress([]byte(to[index])) == false {
				fmt.Printf("地址无效")
				printUsage()
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

	if printChainCmd.Parsed() {


		cli.printchain()
	}

	if addresslistsCmd.Parsed() {


		cli.addressLists()
	}


	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet()
	}

	if createBlockchainCmd.Parsed() {

		if IsValidForAddress([]byte(*flagCreateBlockchainWithAddress)) == false {

			fmt.Println("地址无效")
			printUsage()
			os.Exit(1)
		}


		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {


		if IsValidForAddress([]byte(*getbalanceWithAdress)) == false {

			fmt.Println("地址无效")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getbalanceWithAdress)
	}

}