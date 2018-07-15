package BLC

import "fmt"

// 先用它去查询余额
func (cli *CLI) getBalance(address string,nodeID string)  {

	fmt.Println("地址：" + address)

	// 获取某一个节点的blockchain对象
	blockchain := BlockchainObject(nodeID)
	defer blockchain.LHJ_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.LHJ_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)

}
