package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"os"
)

// 定义数据库名字，常量
const DBNAME  = "BLC.db"

// 表的名字
const BLCTABLENAME  = "BLC"

type Blockchain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}

	return true
}
/*增加区块数据到数据库
 */
func (blc *Blockchain) AddBlockToBlockchain(data string)  {

	err := blc.DB.Update(func(tx *bolt.Tx) error{

		//根据表面获取表信息
		b := tx.Bucket([]byte(BLCTABLENAME))

		if b != nil {

			// 根据最新区块的hash获取最新区块数据
			blockBytes := b.Get(blc.Tip)
			// 将其反序列化
			block := DeserializeBlock(blockBytes)

			//设置要新增的区块的交易数据，高度以及上一区块的hash
			newBlock := NewBlock(data,block.Height + 1,block.Hash)
			//挖矿成功后，将数据序列化后保存到数据库
			err := b.Put(newBlock.Hash,newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//同同时更新最新数据的hash值
			err = b.Put([]byte("l"),newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//修改blockchain的Tip为最新挖矿成功的hash
			blc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}


//创建带有创世区块的区块链
/**
创建创世区块并且初始化数据库信息（创建或者打开数据库）
 */
func CreateBlockchainWithGenesisBlock(data string) *Blockchain {

	// 创建或者打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error{

		//  获取表
		b := tx.Bucket([]byte(BLCTABLENAME))

		//判断表是否存在
		if b == nil {
			// 不在则创建数据库表
			b,err = tx.CreateBucket([]byte(BLCTABLENAME))

			if err != nil {
				log.Panic(err)
			}
		}

		//表存在则设置数据到数据库表中
		if b != nil {
			// 创建创世区块，挖矿
			genesisBlock := CreateGenesisBlock(data)
			// 挖矿，将创世区块存储到表中
			err := b.Put(genesisBlock.Hash,genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 把挖矿成功后的区块信息的hash保存起来
			err = b.Put([]byte("l"),genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			blockHash = genesisBlock.Hash
		}

		return nil
	})

	// 返回区块链对象
	return &Blockchain{blockHash,db}
}
// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain()  {

	var block *Block
	//当前区块的hash
	var cHash []byte = blc.Tip
	for  {
		err := blc.DB.View(func(tx *bolt.Tx) error{

			//根据表名查询表
			b := tx.Bucket([]byte(BLCTABLENAME))
			if b != nil {
				// 根据当前查询出来最新的hash值去查询数据库中的区块细腻
				blockBytes := b.Get(cHash)
				//反序列化表数据
				block = DeserializeBlock(blockBytes)

				fmt.Printf("区块高度:Height：%d\n",block.Height)
				fmt.Printf("上一区块hash:PrevBlockHash：%x\n",block.PrevBlockHash)
				fmt.Printf("交易数据：Data：%s\n",block.Data)
				fmt.Printf("时间戳：Timestamp：%d\n",block.Timestamp)
				fmt.Printf("当前区块：Hash：%x\n",block.Hash)
				fmt.Printf("Nonce：%d\n",block.Nonce)
				fmt.Println()
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		//如果是创世区块了则退出循环
		if big.NewInt(0).Cmp(&hashInt) == 0{
			break;
		}
		cHash = block.PrevBlockHash
	}



}
// 返回Blockchain区块链对象
func BlockchainObject() *Blockchain {

	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BLCTABLENAME))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))
		}


		return nil
	})

	return &Blockchain{tip,db}
}