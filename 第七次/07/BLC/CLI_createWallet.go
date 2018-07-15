package BLC

import "fmt"

func (cli *CLI) LHJ_createWallet(nodeID string)  {

	wallets,_ := NewWallets(nodeID)

	wallets.CreateNewWallet(nodeID)

	fmt.Println(len(wallets.LHJ_WalletsMap))
}
