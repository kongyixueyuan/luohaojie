package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type LHJ_Block struct {
	//1. 区块高度
	LHJ_Height int64
	//2. 上一个区块HASH
	LHJ_PrevBlockHash []byte
	//3. 交易数据
	LHJ_Txs []*LHJ_Transaction
	//4. 时间戳
	LHJ_Timestamp int64
	//5. Hash
	LHJ_Hash []byte
	// 6. Nonce
	LHJ_Nonce int64
}


// 需要将Txs转换成[]byte
func (block *LHJ_Block) LHJ_HashTransactions() []byte  {


	//var txHashes [][]byte
	//var txHash [32]byte
	//
	//for _, tx := range block.Txs {
	//	txHashes = append(txHashes, tx.TxHash)
	//}
	//txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	//
	//return txHash[:]

	var transactions [][]byte

	for _, tx := range block.LHJ_Txs {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data

}


// 将区块序列化成字节数组
func (block *LHJ_Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *LHJ_Block {

	var block LHJ_Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}


//1. 创建新的区块
func LHJ_NewBlock(txs []*LHJ_Transaction,height int64,prevBlockHash []byte) *LHJ_Block {

	//创建区块
	block := &LHJ_Block{height,prevBlockHash,txs,time.Now().Unix(),nil,0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	// 挖矿验证
	hash,nonce := pow.Run()

	block.LHJ_Hash = hash[:]
	block.LHJ_Nonce = nonce

	fmt.Println()

	return block

}

//2. 单独写一个方法，生成创世区块

func LHJ_CreateGenesisBlock(txs []*LHJ_Transaction) *LHJ_Block {


	return LHJ_NewBlock(txs,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

