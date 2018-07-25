package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type LHJ_CLI struct{}

func LHJ_printUsage() {

	fmt.Println()
	fmt.Println("Cli-Usage使用说明 by罗浩杰:")

	fmt.Println("---------------------------------------------------------")
	fmt.Println("\taddresslists -- 显示所有钱包地址.")
	fmt.Println("\tcreatewallet -- 创建钱包.")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -mine -- 交易明细.")
	fmt.Println("\tprintchain -- 显示所有区块信息.")
	fmt.Println("\tgetbalance -address -- 输出交易信息.")
	fmt.Println("\tresetUTXO -- 重置.")
	fmt.Println("\tstartnode -miner ADDRESS -- 启动节点服务器，并且指定挖矿奖励的地址.")
	fmt.Println("---------------------------------------------------------")
	fmt.Println()
}

func LHJ_isValidArgs() {
	if len(os.Args) < 2 {
		LHJ_printUsage()
		os.Exit(1)
	}
}

func (LHJ_cli *LHJ_CLI) LHJ_Run() {

	LHJ_isValidArgs()

	//获取节点ID

	// 设置ID
	// export NODE_ID=3000
	// 读取
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID 还未设置! \n")
		fmt.Printf("格式: export NODE_ID=3000 \n")
		os.Exit(1)
	}

	fmt.Println("===========")
	fmt.Printf("NODE_ID:%s\n", nodeID)
	fmt.Println("===========")

	addresslistsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	resetUTXOCMD := flag.NewFlagSet("resetUTXO", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址.")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址.")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额.")
	flagMine := sendBlockCmd.Bool("mine", false, "是否在当前节点中立即验证.")

	flagMiner := startNodeCmd.String("miner", "", "定义挖矿奖励的地址.")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址.")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.")

	switch os.Args[1] {

	case "addresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "resetUTXO":
		err := resetUTXOCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		LHJ_printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			LHJ_printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)

		for index, fromAddress := range from {
			if LHJ_IsValidForAddress([]byte(fromAddress)) == false || LHJ_IsValidForAddress([]byte(to[index])) == false {
				fmt.Printf("地址无效......")
				LHJ_printUsage()
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flagAmount)
		LHJ_cli.LHJ_send(from, to, amount, nodeID, *flagMine)
	}

	if printChainCmd.Parsed() {

		LHJ_cli.LHJ_printchain(nodeID)
	}

	if resetUTXOCMD.Parsed() {

		fmt.Println("重置UTXO表单......")
		LHJ_cli.LHJ_resetUTXOSet(nodeID)
	}

	if addresslistsCmd.Parsed() {

		LHJ_cli.LHJ_addressLists(nodeID)
	}

	if createWalletCmd.Parsed() {

		// 创建钱包
		LHJ_cli.LHJ_createWallet(nodeID)
	}

	if createBlockchainCmd.Parsed() {

		if LHJ_IsValidForAddress([]byte(*flagCreateBlockchainWithAddress)) == false {
			fmt.Println("地址无效....")
			LHJ_printUsage()
			os.Exit(1)
		}

		LHJ_cli.LHJ_createGenesisBlockchain(*flagCreateBlockchainWithAddress, nodeID)
	}

	if getbalanceCmd.Parsed() {

		if LHJ_IsValidForAddress([]byte(*getbalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			LHJ_printUsage()
			os.Exit(1)
		}

		LHJ_cli.LHJ_getBalance(*getbalanceWithAdress, nodeID)
	}

	if startNodeCmd.Parsed() {

		LHJ_cli.LHJ_startNode(nodeID, *flagMiner)
	}

}
