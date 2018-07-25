package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"fmt"
	"encoding/hex"
	"github.com/boltdb/bolt"
	"os"
)

func LHJ_handleVersion(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {

	var LHJ_buff bytes.Buffer
	var LHJ_payload Version

	dataBytes := LHJ_request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	//Version
	//1. Version
	//2. BestHeight
	//3. 节点地址

	//获取最新Block高度
	bestHeight := LHJ_bc.LHJ_GetBestHeight() //3 1
	//获取接收道德Block高度
	foreignerBestHeight := LHJ_payload.LHJ_BestHeight // 1 3

	if bestHeight > foreignerBestHeight {
		LHJ_sendVersion(LHJ_payload.LHJ_AddrFrom,LHJ_bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		LHJ_sendGetBlocks(LHJ_payload.LHJ_AddrFrom)
	}

	if !LHJ_nodeIsKnown(LHJ_payload.LHJ_AddrFrom) {
		LHJ_knowNodes = append(LHJ_knowNodes, LHJ_payload.LHJ_AddrFrom)
	}

}

func LHJ_handleAddr(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {



}

func LHJ_handleGetblocks(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {


	var LHJ_buff bytes.Buffer
	var LHJ_payload LHJ_GetBlocks

	dataBytes := LHJ_request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := LHJ_bc.LHJ_GetBlockHashes()

	//txHash blockHash
	LHJ_sendInv(LHJ_payload.LHJ_AddrFrom, BLOCK_TYPE, blocks)


}

func LHJ_handleGetData(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {

	var LHJ_buff bytes.Buffer
	var LHJ_payload LHJ_GetData

	dataBytes := LHJ_request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	if LHJ_payload.LHJ_Type == BLOCK_TYPE {

		block, err := LHJ_bc.LHJ_GetBlock([]byte(LHJ_payload.LHJ_Hash))
		if err != nil {
			return
		}

		LHJ_sendBlock(LHJ_payload.LHJ_AddrFrom, block)
	}

	if LHJ_payload.LHJ_Type == TX_TYPE {

		tx := LHJ_memoryTxPool[hex.EncodeToString(LHJ_payload.LHJ_Hash)]

		LHJ_sendTx(LHJ_payload.LHJ_AddrFrom,tx)

	}
}

func LHJ_handleBlock(request []byte,LHJ_bc *LHJ_Blockchain)  {
	var LHJ_buff bytes.Buffer
	var LHJ_payload LHJ_BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := LHJ_payload.LHJ_Block

	block := DeserializeBlock(blockBytes)

	fmt.Println("Recevied a new block!")
	LHJ_bc.LHJ_AddBlock(block)

	fmt.Printf("Added block %x\n", block.LHJ_Hash)

	if len(LHJ_transactionArray) > 0 {
		blockHash := LHJ_transactionArray[0]
		LHJ_sendGetData(LHJ_payload.LHJ_AddrFrom, "block", blockHash)

		LHJ_transactionArray = LHJ_transactionArray[1:]
	} else {

		fmt.Println("数据库重置......")
		UTXOSet := &LHJ_UTXOSet{LHJ_bc}
		UTXOSet.LHJ_ResetUTXOSet()

	}

}

func LHJ_handleTx(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {

	var LHJ_buff bytes.Buffer
	var LHJ_payload Tx

	dataBytes := LHJ_request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	tx := LHJ_payload.LHJ_Tx
	LHJ_memoryTxPool[hex.EncodeToString(tx.LHJ_TxHash)] = tx

	// 说明主节点自己
	if LHJ_nodeAddress == LHJ_knowNodes[0] {
		// 给矿工节点发送交易hash
		for _,nodeAddr := range LHJ_knowNodes {

			if nodeAddr != LHJ_nodeAddress && nodeAddr != LHJ_payload.LHJ_AddrFrom {
				LHJ_sendInv(nodeAddr,TX_TYPE,[][]byte{tx.LHJ_TxHash})
			}

		}
	}

	// 矿工进行挖矿验证
	if len(LHJ_minerAddress) > 0 {

		LHJ_bc.LHJ_DB.Close()

		blockchain := LHJ_BlockchainObject(os.Getenv("NODE_ID"))
		//

		defer blockchain.LHJ_DB.Close()

		utxoSet := &LHJ_UTXOSet{LHJ_bc}

		var txs []*LHJ_Transaction

		txs = append(txs, tx)

		//奖励
		coinbaseTx := LHJ_NewCoinbaseTransaction(LHJ_minerAddress)
		txs = append(txs,coinbaseTx)



		_txs := []*LHJ_Transaction{}

		fmt.Println("开始进行数字签名验证.....")

		for _,tx := range txs  {

			// 作业，数字签名失败
			if LHJ_bc.LHJ_VerifyTransaction(tx,_txs) != true {
				log.Panic("ERROR: Invalid transaction")
			}

			_txs = append(_txs,tx)
		}

		fmt.Println("数字签名验证成功.....")

		//1. 通过相关算法建立Transaction数组
		var LHJ_block *LHJ_Block

		LHJ_bc.LHJ_DB.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(LHJ_blockTableName))
			if b != nil {

				hash := b.Get([]byte("l"))

				blockBytes := b.Get(hash)

				LHJ_block = DeserializeBlock(blockBytes)

			}

			return nil
		})

		//2. 建立新的区块
		LHJ_block = LHJ_NewBlock(txs, LHJ_block.LHJ_Height+1, LHJ_block.LHJ_Hash)

		//将新区块存储到数据库
		LHJ_bc.LHJ_DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(LHJ_blockTableName))
			if b != nil {

				b.Put(LHJ_block.LHJ_Hash, LHJ_block.Serialize())

				b.Put([]byte("l"), LHJ_block.LHJ_Hash)

				LHJ_bc.LHJ_Tip = LHJ_block.LHJ_Hash

			}
			return nil
		})
		utxoSet.LHJ_Update()
		LHJ_sendBlock(LHJ_knowNodes[0],LHJ_block.Serialize())
	}


}


func LHJ_handleInv(LHJ_request []byte,LHJ_bc *LHJ_Blockchain)  {

	var LHJ_buff bytes.Buffer
	var LHJ_payload Inv

	dataBytes := LHJ_request[COMMANDLENGTH:]

	// 反序列化
	LHJ_buff.Write(dataBytes)
	dec := gob.NewDecoder(&LHJ_buff)
	err := dec.Decode(&LHJ_payload)
	if err != nil {
		log.Panic(err)
	}

	// Ivn 3000 block hashes [][]

	if LHJ_payload.LHJ_Type == BLOCK_TYPE {

		//tansactionArray = payload.Items

		//payload.Items

		blockHash := LHJ_payload.LHJ_Items[0]
		LHJ_sendGetData(LHJ_payload.LHJ_AddrFrom, BLOCK_TYPE , blockHash)

		if len(LHJ_payload.LHJ_Items) >= 1 {
			LHJ_transactionArray = LHJ_payload.LHJ_Items[1:]
		}
	}

	if LHJ_payload.LHJ_Type == TX_TYPE {

		txHash := LHJ_payload.LHJ_Items[0]
		if LHJ_memoryTxPool[hex.EncodeToString(txHash)] == nil  {
			LHJ_sendGetData(LHJ_payload.LHJ_AddrFrom, TX_TYPE , txHash)
		}

	}

}