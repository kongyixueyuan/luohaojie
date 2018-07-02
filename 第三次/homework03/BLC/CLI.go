
package BLC
import (
"fmt"
"os"
"flag"
"log"
)

type CLI struct {
	BC *Blockchain
}



func printUsage()  {

	fmt.Println("使用方法:")
	fmt.Println("\tcreateGenesisBlockchain -data -- 创世区块数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")

}

func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string)  {
	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain(data)
}

/*
查询，迭代输出区块链信息
 */
func (cli *CLI) printchain()  {
	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	//延迟关闭数据库，程序运行到最后执行
	defer blockchain.DB.Close()

	blockchain.Printchain()
}

func (cli *CLI) createGenesisBlockchain(data string)  {
/*
创建创世区块并且初始化数据库信息（创建或者打开数据库）
 */
	blockchain := CreateBlockchainWithGenesisBlock(data)
	//延迟关闭数据库，程序运行到最后执行
	defer blockchain.DB.Close()


}


func (cli *CLI) Run()  {

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createGenesisBlockchainCmd := flag.NewFlagSet("createGenesisBlockchain",flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","","交易数据......")

	flagCreateGenesisBlockchainWithData := createGenesisBlockchainCmd.String("data","创世区块的交易数据","这是创建创世区块的命令")


	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createGenesisBlockchain":
		err := createGenesisBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createGenesisBlockchainCmd.Parsed() {

		if *flagCreateGenesisBlockchainWithData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateGenesisBlockchainWithData)
	}

}
