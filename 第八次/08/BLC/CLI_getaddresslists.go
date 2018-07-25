package BLC

import "fmt"

//打印所有钱包地址
func (LHJ_cli *LHJ_CLI) LHJ_addressLists(LHJ_nodeID string) {

	fmt.Println("======================================")

	fmt.Println("打印所有钱包地址: ")

	wallets, _ := LHJ_NewWallets(LHJ_nodeID)

	for address, _ := range wallets.LHJ_WalletsMap {

		fmt.Println(address)

	}

	fmt.Println("======================================")
}
