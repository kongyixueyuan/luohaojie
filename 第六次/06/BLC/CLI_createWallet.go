package BLC

import "fmt"

func (cli *CLI) createWallet()  {

	wallets,_ := LHJ_NewWallets()

	wallets.LHJ_CreateNewWallet()

	fmt.Println(len(wallets.LHJ_WalletsMap))
}
