package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printchain()  {

	//判断DB是否存在
	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()

}