package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type LHJ_BlockchainIterator struct {
	LHJ_CurrentHash []byte
	LHJ_DB  *bolt.DB
}

func (blockchainIterator *LHJ_BlockchainIterator) Next() *LHJ_Block {

	var block *LHJ_Block

	err := blockchainIterator.LHJ_DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.LHJ_CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.LHJ_CurrentHash = block.LHJ_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}


	return block

}