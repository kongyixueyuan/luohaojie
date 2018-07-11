package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printchain()  {

	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := LHJ_BlockchainObject()

	defer blockchain.LHJ_DB.Close()

	blockchain.LHJ_Printchain()

}