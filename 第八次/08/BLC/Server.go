package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

func LHJ_startServer(LHJ_nodeID string, LHJ_minerAdd string) {

	// 当前节点的IP地址
	LHJ_nodeAddress = fmt.Sprintf("localhost:%s", LHJ_nodeID)

	LHJ_minerAddress = LHJ_minerAdd

	ln, err := net.Listen(PROTOCOL, LHJ_nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	bc := LHJ_BlockchainObject(LHJ_nodeID)

	defer bc.LHJ_DB.Close()

	// 第一个终端：端口为3000,启动的就是主节点
	// 第二个终端：端口为3001，钱包节点
	// 第三个终端：端口号为3002，矿工节点
	if LHJ_nodeAddress != LHJ_knowNodes[0] {

		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		LHJ_sendVersion(LHJ_knowNodes[0], bc)

	}

	for {

		// 收到的数据的格式是固定的，12字节+结构体字节数组
		// 接收客户端发送过来的数据
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go LHJ_handleConnection(conn, bc)

	}

}

func LHJ_handleConnection(LHJ_conn net.Conn, LHJ_bc *LHJ_Blockchain) {

	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(LHJ_conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n", request[:COMMANDLENGTH])

	//version
	command := bytesToCommand(request[:COMMANDLENGTH])

	// 12字节 + 某个结构体序列化以后的字节数组

	switch command {
	case COMMAND_VERSION:
		LHJ_handleVersion(request, LHJ_bc)
	case COMMAND_ADDR:
		LHJ_handleAddr(request, LHJ_bc)
	case COMMAND_BLOCK:
		LHJ_handleBlock(request, LHJ_bc)
	case COMMAND_GETBLOCKS:
		LHJ_handleGetblocks(request, LHJ_bc)
	case COMMAND_GETDATA:
		LHJ_handleGetData(request, LHJ_bc)
	case COMMAND_INV:
		LHJ_handleInv(request, LHJ_bc)
	case COMMAND_TX:
		LHJ_handleTx(request, LHJ_bc)
	default:
		fmt.Println("Unknown command!")
	}

	LHJ_conn.Close()
}

func LHJ_nodeIsKnown(addr string) bool {
	for _, node := range LHJ_knowNodes {
		if node == addr {
			return true
		}
	}

	return false
}
