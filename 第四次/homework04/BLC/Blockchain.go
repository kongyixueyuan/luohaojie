package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"strconv"
	"encoding/hex"
)

// 定义数据库名字，常量
const DBNAME  = "BLC.db"

// 表的名字
const BLCTABLENAME  = "BLC"

type Blockchain struct {
	Tip []byte //保存最新的区块的Hash
	DB  *bolt.DB
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}

	return true
}





// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("区块高度:Height：%d\n",block.Height)
		fmt.Printf("上一区块hash:PrevBlockHash：%x\n",block.PrevBlockHash)
		//格式化时间输出
		fmt.Printf("时间戳：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("当前区块：Hash：%x\n",block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}
		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		//如果是创世区块了则退出循环
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}

/*增加区块数据到数据库
 */
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction)  {

	err := blc.DB.Update(func(tx *bolt.Tx) error{

		//根据表面获取表信息
		b := tx.Bucket([]byte(BLCTABLENAME))

		if b != nil {

			// 根据最新区块的hash获取最新区块数据
			blockBytes := b.Get(blc.Tip)
			// 将其反序列化
			block := DeserializeBlock(blockBytes)

			//设置要新增的区块的交易数据，高度以及上一区块的hash
			newBlock := NewBlock(txs,block.Height + 1,block.Hash)
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
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {

	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已创建")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块")

	// 创建或者打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	// 关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(BLCTABLENAME))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)

			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.Hash
		}

		return nil
	})

	return &Blockchain{genesisHash, db}

}

// 返回Blockchain对象
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

	return &Blockchain{tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string,txs []*Transaction) []*UTXO {



	var unUTXOs []*UTXO

	//存储txhash:outputIndex下标数组
	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {
			//是否创世区块
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockWithAddress(address) {

					key := hex.EncodeToString(in.TxHash)
					//把所有未花费的txoutput 追加到map中
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}


	for _,tx := range txs {

		w1:
		for index,out := range tx.Vouts {

			if out.UnLockScriptPubKeyWithAddress(address) {
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue w1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}




	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.Next()

		for i := len(block.Txs) - 1; i >= 0 ; i-- {

			tx := block.Txs[i]

			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {

						key := hex.EncodeToString(in.TxHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		w2:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {


					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue w2
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//如果是创世区块了则退出循环
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int,txs []*Transaction) (int64, map[string][]int) {

	//获取所有的UTXO

	utxos := blockchain.UnUTXOs(from,txs)

//存储txhash:outputIndex下标数组
	spendableUTXO := make(map[string][]int)

	// 遍历utxos

	//记录金额
	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		//迭代++,如果满足转账金额，则跳出循环
		if value >= int64(amount) {
			break
		}
	}

	//如果迭代完都不够钱，则提示，退出
	if value < int64(amount) {

		fmt.Printf("%s's fund is not enough\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	// 通过相关算法建立Transantion数组

	//建立一笔交易

	var txs []*Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain,txs)
		txs = append(txs, tx)
	}


	//通过相关算法建立Transaction数组
	var block *Block

	blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(BLCTABLENAME))
		if b != nil {

			hash := b.Get([]byte("l"))


			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)
			fmt.Println(block)

		}

		return nil
	})

	//建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)
	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLCTABLENAME))
		if b != nil {

			b.Put(block.Hash, block.Serialize())

			b.Put([]byte("l"), block.Hash)

			blockchain.Tip = block.Hash

		}
		return nil
	})

}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {

	utxos := blockchain.UnUTXOs(address,[]*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Output.Value
	}

	return amount
}
