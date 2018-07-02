package BLC

import (
	"time"
	"fmt"
)

type Block struct {
	Height     int64  //区块编号
	TxData     []byte //交易数据
	PHash      []byte //上一区块hash
	CreateTime int64  //区块创建时间戳
	Hash       []byte //当前区块hash
	Nonce      int64
}

func NewBlock(height int64, txData string, pHash []byte) *Block {
	block := &Block{height, []byte(txData), pHash, time.Now().Unix(), nil, 0}
	//fmt.Println(block)

	pow := NewPow(block)
	// 挖矿验证
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Println()

	return block

}
/*
创世区块，调用NewBlock
 */
func GenesisBlock(data string) *Block {
	return NewBlock(1,data, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
