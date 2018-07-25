package BLC

import "fmt"

func (LHJ_cli *LHJ_CLI) LHJ_createWallet(LHJ_nodeId string) {

	wallets,_ := LHJ_NewWallets(LHJ_nodeId)

	wallets.LHJ_CreateNewWallet(LHJ_nodeId)

	fmt.Println(len(wallets.LHJ_WalletsMap))

}