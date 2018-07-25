package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type LHJ_BlockchainIterator struct {
	LHJ_CurrentHash []byte
	LHJ_DB          *bolt.DB
}

func (LHJ_blockchainIterator *LHJ_BlockchainIterator) LHJ_Next() *LHJ_Block {
	var LHJ_block *LHJ_Block

	err := LHJ_blockchainIterator.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {
			LHJ_currentBlockBytes := b.Get(LHJ_blockchainIterator.LHJ_CurrentHash)
			// 获取到当前迭代器里面的currentHash所对应的区块
			LHJ_block = DeserializeBlock(LHJ_currentBlockBytes)
			// 更新迭代器里面的currentHash
			LHJ_blockchainIterator.LHJ_CurrentHash = LHJ_block.LHJ_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return LHJ_block
}
