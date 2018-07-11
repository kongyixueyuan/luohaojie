package BLC

import (
	"fmt"
	"os"
)

// 转账
func (cli *CLI) send(from []string,to []string,amount []string)  {


	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := LHJ_BlockchainObject()
	defer blockchain.LHJ_DB.Close()

	blockchain.LHJ_MineNewBlock(from,to,amount)

	utxoSet := &UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.Update()

}

