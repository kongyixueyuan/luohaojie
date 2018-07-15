package BLC

import (
	"fmt"
	"io"
	"bytes"
	"log"
	"net"
)


//COMMAND_VERSION
func sendVersion(toAddress string,bc *LHJ_Blockchain)  {


	bestHeight := bc.LHJ_GetBestHeight()
	payload := gobEncode(Version{LHJ_NODE_VERSION, bestHeight, nodeAddress})

	request := append(commandToBytes(LHJ_COMMAND_VERSION), payload...)

	sendData(toAddress,request)


}



//COMMAND_GETBLOCKS
func sendGetBlocks(toAddress string)  {

	payload := gobEncode(GetBlocks{nodeAddress})

	request := append(commandToBytes(LHJ_COMMAND_GETBLOCKS), payload...)

	sendData(toAddress,request)

}

// 主节点将自己的所有的区块hash发送给钱包节点
//COMMAND_BLOCK
//
func sendInv(toAddress string, kind string, hashes [][]byte) {

	payload := gobEncode(Inv{nodeAddress,kind,hashes})

	request := append(commandToBytes(LHJ_COMMAND_INV), payload...)

	sendData(toAddress,request)

}



func sendGetData(toAddress string, kind string ,blockHash []byte) {

	payload := gobEncode(GetData{nodeAddress,kind,blockHash})

	request := append(commandToBytes(LHJ_COMMAND_GETDATA), payload...)

	sendData(toAddress,request)
}


func sendData(to string,data []byte)  {

	fmt.Println("客户端向服务器发送数据......")
	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}