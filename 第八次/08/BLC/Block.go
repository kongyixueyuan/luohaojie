package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"fmt"
)

type LHJ_Block struct {
	LHJ_Height        int64
	LHJ_PrevBlockHash []byte
	LHJ_Txs           []*LHJ_Transaction
	LHJ_Timestamp     int64
	LHJ_Hash          []byte
	LHJ_Nonce         int64
}

//序列化成字节数组
func (block *LHJ_Block) LHJ_HashTransactions() []byte {

	var transactions [][]byte

	for _, tx := range block.LHJ_Txs {
		transactions = append(transactions, tx.LHJ_Serialize())
	}
	mTree := LHJ_NewMerkleTree(transactions)

	return mTree.LHJ_RootNode.LHJ_Data
}

//序列化
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

//创建创世区块
func LHJ_NewBlock(LHJ_txs []*LHJ_Transaction, height int64, prevBlockHash []byte) *LHJ_Block {

	//创建区块
	block := &LHJ_Block{height, prevBlockHash, LHJ_txs, time.Now().Unix(), nil, 0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := LHJ_NewProofOfWork(block)

	// 挖矿验证
	hash, nonce := pow.LHJ_Run()

	block.LHJ_Hash = hash[:]
	block.LHJ_Nonce = nonce

	fmt.Println()

	return block
}

//生成创世区块
func LHJ_CreateGenesisBlock(txs []*LHJ_Transaction) *LHJ_Block {
	return LHJ_NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
