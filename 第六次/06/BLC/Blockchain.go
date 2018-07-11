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
	"crypto/ecdsa"
	"bytes"
)

// 数据库名字
const LHJ_dbName = "LHJ_blockchain.db"

// 表的名字
const LHJ_blockTableName = "LHJ_blocks"

type LHJ_Blockchain struct {
	LHJ_Tip []byte //最新的区块的Hash
	LHJ_DB  *bolt.DB
}

// 迭代器
func (blockchain *LHJ_Blockchain) Iterator() *LHJ_BlockchainIterator {

	return &LHJ_BlockchainIterator{blockchain.LHJ_Tip, blockchain.LHJ_DB}
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(LHJ_dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 遍历输出所有区块的信息
func (blc *LHJ_Blockchain) LHJ_Printchain() {

	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.LHJ_Height)
		fmt.Printf("PrevBlockHash：%x\n", block.LHJ_PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.LHJ_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.LHJ_Hash)
		fmt.Printf("Nonce：%d\n", block.LHJ_Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.LHJ_Txs {

			fmt.Printf("%x\n", tx.LHJ_TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.LHJ_Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%x\n", in.PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.LHJ_Vouts {
				//fmt.Println(out.Value)
				fmt.Printf("%d\n",out.Value)
				//fmt.Println(out.Ripemd160Hash)
				fmt.Printf("%x\n",out.Ripemd160Hash)
			}
		}

		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}

//// 增加区块到区块链里面
func (blc *LHJ_Blockchain) LHJ_AddBlockToBlockchain(txs []*LHJ_Transaction) {

	err := blc.LHJ_DB.Update(func(tx *bolt.Tx) error {

		//1. 获取表
		b := tx.Bucket([]byte(LHJ_blockTableName))
		//2. 创建新区块
		if b != nil {

			// ⚠️，先获取最新区块
			blockBytes := b.Get(blc.LHJ_Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)

			//3. 将区块序列化并且存储到数据库中
			newBlock := LHJ_NewBlock(txs, block.LHJ_Height+1, block.LHJ_Hash)
			err := b.Put(newBlock.LHJ_Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), newBlock.LHJ_Hash)
			if err != nil {
				log.Panic(err)
			}
			//5. 更新blockchain的Tip
			blc.LHJ_Tip = newBlock.LHJ_Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//1. 创建带有创世区块的区块链
func LHJ_CreateBlockchainWithGenesisBlock(address string) *LHJ_Blockchain {

	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(LHJ_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	// 关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(LHJ_blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)

			genesisBlock := LHJ_CreateGenesisBlock([]*LHJ_Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.LHJ_Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.LHJ_Hash)
			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.LHJ_Hash
		}

		return nil
	})

	return &LHJ_Blockchain{genesisHash, db}

}

// 返回Blockchain对象
func LHJ_BlockchainObject() *LHJ_Blockchain {

	db, err := bolt.Open(LHJ_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))

		}

		return nil
	})

	return &LHJ_Blockchain{tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *LHJ_Blockchain) LHJ_UnUTXOs(address string,txs []*LHJ_Transaction) []*UTXO {



	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.LHJ_Vins {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(address))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
				if in.UnLockRipemd160Hash(ripemd160Hash) {

					key := hex.EncodeToString(in.TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}


	for _,tx := range txs {

		Work1:
		for index,out := range tx.LHJ_Vouts {

			if out.UnLockScriptPubKeyWithAddress(address) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(address)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.LHJ_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.LHJ_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.LHJ_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.LHJ_TxHash, index, out}
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

		fmt.Println(block)
		fmt.Println()

		for i := len(block.LHJ_Txs) - 1; i >= 0 ; i-- {

			tx := block.LHJ_Txs[i]
			// txHash
			// Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.LHJ_Vins {
					//是否能够解锁
					publicKeyHash := Base58Decode([]byte(address))

					ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

					if in.UnLockRipemd160Hash(ripemd160Hash) {

						key := hex.EncodeToString(in.TxHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.LHJ_Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {

					fmt.Println(out)
					fmt.Println(spentTXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.LHJ_TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.LHJ_TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.LHJ_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *LHJ_Blockchain) LHJ_FindSpendableUTXOS(from string, amount int,txs []*LHJ_Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.LHJ_UnUTXOs(from,txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {

		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *LHJ_Blockchain) LHJ_MineNewBlock(from []string, to []string, amount []string) {

	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易


	utxoSet := &UTXOSet{blockchain}

	var txs []*LHJ_Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//奖励
	tx := NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)


	//1. 通过相关算法建立Transaction数组
	var block *LHJ_Block

	blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)

		}

		return nil
	})


	// 在建立新区块之前对txs进行签名验证

	_txs := []*LHJ_Transaction{}

	for _,tx := range txs  {

		if blockchain.LHJ_VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}


	//2. 建立新的区块
	block = LHJ_NewBlock(txs, block.LHJ_Height+1, block.LHJ_Hash)

	//将新区块存储到数据库
	blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LHJ_blockTableName))
		if b != nil {

			b.Put(block.LHJ_Hash, block.Serialize())

			b.Put([]byte("l"), block.LHJ_Hash)

			blockchain.LHJ_Tip = block.LHJ_Hash

		}
		return nil
	})

}

// 查询余额
func (blockchain *LHJ_Blockchain) LHJ_GetBalance(address string) int64 {

	utxos := blockchain.LHJ_UnUTXOs(address,[]*LHJ_Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Output.Value
	}

	return amount
}

func (bclockchain *LHJ_Blockchain) LHJ_SignTransaction(tx *LHJ_Transaction,privKey ecdsa.PrivateKey,txs []*LHJ_Transaction)  {

	if tx.IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]LHJ_Transaction)

	for _, vin := range tx.LHJ_Vins {
		prevTX, err := bclockchain.LHJ_FindTransaction(vin.TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.LHJ_TxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)

}


func (bc *LHJ_Blockchain) LHJ_FindTransaction(ID []byte,txs []*LHJ_Transaction) (LHJ_Transaction, error) {


	for _,tx := range txs  {
		if bytes.Compare(tx.LHJ_TxHash, ID) == 0 {
			return *tx, nil
		}
	}


	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.LHJ_Txs {
			if bytes.Compare(tx.LHJ_TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return LHJ_Transaction{},nil
}


// 验证数字签名
func (bc *LHJ_Blockchain) LHJ_VerifyTransaction(tx *LHJ_Transaction,txs []*LHJ_Transaction) bool {


	prevTXs := make(map[string]LHJ_Transaction)

	for _, vin := range tx.LHJ_Vins {
		prevTX, err := bc.LHJ_FindTransaction(vin.TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.LHJ_TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}


// [string]*TXOutputs
func (blc *LHJ_Blockchain) LHJ_FindUTXOMap() map[string]*TXOutputs  {

	blcIterator := blc.Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)


	utxoMaps := make(map[string]*TXOutputs)


	for {
		block := blcIterator.Next()


		for i := len(block.LHJ_Txs) - 1; i >= 0 ;i-- {

			txOutputs := &TXOutputs{[]*UTXO{}}

			tx := block.LHJ_Txs[i]


			// coinbase
			if tx.IsCoinbaseTransaction() == false {
				for _,txInput := range tx.LHJ_Vins {

					txHash := hex.EncodeToString(txInput.TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}



			txHash := hex.EncodeToString(tx.LHJ_TxHash)

			WorkOutLoop:
			for index,out := range tx.LHJ_Vouts  {

				if tx.IsCoinbaseTransaction() {

					fmt.Println("IsCoinbaseTransaction")
					fmt.Println(out)
					fmt.Println(txHash)
				}

				txInputs := spentableUTXOsMap[txHash]

				if len(txInputs) > 0 {

					isSpent := false

					for _,in := range  txInputs {

						outPublicKey := out.Ripemd160Hash
						inPublicKey := in.PublicKey

						if bytes.Compare(outPublicKey,LHJ_Ripemd160Hash(inPublicKey)) == 0{
							if index == in.Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}

					}

					if isSpent == false {
						utxo := &UTXO{tx.LHJ_TxHash,index,out}
						txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
					}

				} else {
					utxo := &UTXO{tx.LHJ_TxHash,index,out}
					txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
				}

			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}



	}

	return utxoMaps
}