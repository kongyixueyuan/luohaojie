package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI) addressLists()  {

	fmt.Println("输出所有的钱包地址:")

	wallets,_ := NewWallets()

	for address,_ := range wallets.Wallets {

		fmt.Println(address)
	}
}