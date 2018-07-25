package BLC

import (
	"fmt"
	"os"
)

func (LHJ_cli *LHJ_CLI) LHJ_startNode(LHJ_nodeID string,LHJ_minerAdd string)  {

	// 启动服务器
	if LHJ_minerAdd == "" || LHJ_IsValidForAddress([]byte(LHJ_minerAdd))  {

		//  启动服务器
		fmt.Printf("启动服务器:localhost:%s\n",LHJ_nodeID)
		LHJ_startServer(LHJ_nodeID,LHJ_minerAdd)

	} else {

		fmt.Println("指定的地址无效....")
		os.Exit(0)

	}

}