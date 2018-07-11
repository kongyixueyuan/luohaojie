package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI) addressLists()  {

	fmt.Println("打印输出所有的钱包地址:")

	wallets,_ := LHJ_NewWallets()

	for address,_ := range wallets.LHJ_WalletsMap {

		fmt.Println(address)
	}
}