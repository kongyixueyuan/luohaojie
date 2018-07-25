package BLC

import (
	"io"
	"bytes"
	"log"
	"net"
)

//COMMAND_VERSION
func LHJ_sendVersion(LHJ_toAddress string, LHJ_bc *LHJ_Blockchain) {

	//返回当前区块链最新区块高度
	bestHeight := LHJ_bc.LHJ_GetBestHeight()

	//创建version对象,并且序列化
	payload := gobEncode(Version{NODE_VERSION, bestHeight, LHJ_nodeAddress})

	//version
	request := append(commandToBytes(COMMAND_VERSION), payload...)

	LHJ_sendData(LHJ_toAddress, request)

}

//COMMAND_GETBLOCKS
func LHJ_sendGetBlocks(LHJ_toAddress string) {

	payload := gobEncode(LHJ_GetBlocks{LHJ_nodeAddress})

	request := append(commandToBytes(COMMAND_GETBLOCKS), payload...)

	LHJ_sendData(LHJ_toAddress, request)

}

// 主节点将自己的所有的区块hash发送给钱包节点
//COMMAND_BLOCK
func LHJ_sendInv(LHJ_toAddress string, LHJ_kind string, LHJ_hashes [][]byte) {

	payload := gobEncode(Inv{LHJ_nodeAddress, LHJ_kind, LHJ_hashes})

	request := append(commandToBytes(COMMAND_INV), payload...)

	LHJ_sendData(LHJ_toAddress, request)

}

func LHJ_sendGetData(LHJ_toAddress string, LHJ_kind string, LHJ_blockHash []byte) {

	payload := gobEncode(LHJ_GetData{LHJ_nodeAddress, LHJ_kind, LHJ_blockHash})

	request := append(commandToBytes(COMMAND_GETDATA), payload...)

	LHJ_sendData(LHJ_toAddress, request)
}

func LHJ_sendBlock(LHJ_toAddress string, LHJ_block []byte) {

	payload := gobEncode(LHJ_BlockData{LHJ_nodeAddress, LHJ_block})

	request := append(commandToBytes(COMMAND_BLOCK), payload...)

	LHJ_sendData(LHJ_toAddress, request)

}

func LHJ_sendTx(LHJ_toAddress string, LHJ_tx *LHJ_Transaction) {

	payload := gobEncode(Tx{LHJ_nodeAddress, LHJ_tx})

	request := append(commandToBytes(COMMAND_TX), payload...)

	//fmt.Println(LHJ_toAddress)
	//fmt.Println(request)

	LHJ_sendData(LHJ_toAddress, request)

}

func LHJ_sendData(LHJ_to string, LHJ_data []byte) {

	conn, err := net.Dial("tcp", LHJ_to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(LHJ_data))
	if err != nil {
		log.Panic(err)
	}
}
