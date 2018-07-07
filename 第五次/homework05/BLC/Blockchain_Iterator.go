package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器
type BlockchainIterator struct {
	//当前遍历hash
	CurrentHash []byte
	//区块链数据库
	DB *bolt.DB
}

func (blockchainIterator *BlockchainIterator) Next() *Block {

	var block *Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}


	return block

}