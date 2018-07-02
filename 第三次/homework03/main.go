package main

import "homework03/BLC"

/*func main()  {

	// 创建创世区块，并且初始化数据库信息
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	//延迟关闭数据库，程序运行到最后执行
	defer blockchain.DB.Close()
fmt.Println()
	//新建区块信息
	blockchain.AddBlockToBlockchain("罗浩杰交易数据一，购买1000个BTC")
	fmt.Println()

	//便利所有区块信息

	blockchain.Printchain();
	// ./main -printchain
	flagString := flag.String("printchain","","循环遍历所有的区块信息")
	flag.Parse()
	fmt.Printf("%s\n",*flagString)


}*/
func main()  {

	cli := BLC.CLI{}
	cli.Run()
}
