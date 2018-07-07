package BLC

func (cli *CLI) createWallet()  {

	wallets,_ := NewWallets()

	wallets.CreateWallet()

	//fmt.Println(len(wallets.Wallets))
}
