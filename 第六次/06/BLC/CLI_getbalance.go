package BLC

import "fmt"

// 先用它去查询余额
func (cli *CLI) getBalance(address string)  {

	fmt.Println("钱包地址：" + address)

	blockchain := LHJ_BlockchainObject()
	defer blockchain.LHJ_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)

}
