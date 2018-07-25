package BLC

import "fmt"

func (LHJ_cli *LHJ_CLI) LHJ_getBalance(LHJ_address string,LHJ_nodeID string)  {

	fmt.Println()
	fmt.Println("地址：" + LHJ_address)
	fmt.Println()

	// 获取某一个节点的blockchain对象
	blockchain := LHJ_BlockchainObject(LHJ_nodeID)
	defer blockchain.LHJ_DB.Close()

	utxoSet := &LHJ_UTXOSet{blockchain}

	amount := utxoSet.LHJ_GetBalance(LHJ_address)

	fmt.Printf("%s 一共有 %d 个Token\n",LHJ_address,amount)
	fmt.Println()

}