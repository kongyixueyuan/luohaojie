package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"time"
	"math/big"
	"log"
	"encoding/hex"
	"strconv"
	"crypto/ecdsa"
	"bytes"
)

//库名
const LHJ_dbName = "LHJ_blockchain_%s.db"

//表名
const LHJ_blockTableName = "LHJ_blocks"

type LHJ_Blockchain struct {
	LHJ_Tip []byte
	LHJ_DB  *bolt.DB
}

//迭代器

func (blockchain *LHJ_Blockchain) LHJ_Iterator() *LHJ_BlockchainIterator {
	return &LHJ_BlockchainIterator{blockchain.LHJ_Tip, blockchain.LHJ_DB}
}

//判断数据库存在
func LHJ_DBExists(LHJ_dbName string) bool {
	if _, err := os.Stat(LHJ_dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//遍历所有区块信息
func (LHJ_blc *LHJ_Blockchain) LHJ_Printchain() {

	fmt.Println()
	fmt.Println("===========LHJ_遍历所有区块===========")
	fmt.Println()

	LHJ_blockchainIterator := LHJ_blc.LHJ_Iterator()

	for {
		//获取区块
		LHJ_block := LHJ_blockchainIterator.LHJ_Next()

		fmt.Printf("Height：%d\n", LHJ_block.LHJ_Height)
		fmt.Printf("PrevBlockHash：%x\n", LHJ_block.LHJ_PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(LHJ_block.LHJ_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", LHJ_block.LHJ_Hash)
		fmt.Printf("Nonce：%d\n", LHJ_block.LHJ_Nonce)
		fmt.Println("Txs:")

		//遍历交易信息
		for _, tx := range LHJ_block.LHJ_Txs {

			fmt.Printf("%x\n", tx.LHJ_TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.LHJ_Vins {
				fmt.Printf("%x\n", in.LHJ_TxHash)
				fmt.Printf("%d\n", in.LHJ_Vout)
				fmt.Printf("%x\n", in.LHJ_PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.LHJ_Vouts {
				//fmt.Println(out.Value)
				fmt.Printf("%d\n", out.LHJ_Value)
				//fmt.Println(out.Ripemd160Hash)
				fmt.Printf("%x\n", out.LHJ_Ripemd160Hash)
			}
		}

		fmt.Println("=====================================")
		fmt.Println()

		var LHJ_hashInt big.Int
		LHJ_hashInt.SetBytes(LHJ_block.LHJ_PrevBlockHash)

		if big.NewInt(0).Cmp(&LHJ_hashInt) == 0 {
			break
		}
	}
}

// 增加区块到区块链里面
func (LHJ_blc *LHJ_Blockchain) LHJ_AddBlockToBlockchain(LHJ_txs []*LHJ_Transaction) {

	err := LHJ_blc.LHJ_DB.Update(func(tx *bolt.Tx) error {

		//1. 获取表
		b := tx.Bucket([]byte(LHJ_blockTableName))
		//2. 创建新区块
		if b != nil {

			// 先获取最新区块
			blockBytes := b.Get(LHJ_blc.LHJ_Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)

			//3. 将区块序列化并且存储到数据库中
			newBlock := LHJ_NewBlock(LHJ_txs, block.LHJ_Height+1, block.LHJ_Hash)
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
			LHJ_blc.LHJ_Tip = newBlock.LHJ_Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 创建带有创世区块的区块链
func LHJ_CreateBlockchainWithGenesisBlock(LHJ_address string, LHJ_nodeID string) *LHJ_Blockchain {

	// 格式化数据库名字
	LHJ_dbName := fmt.Sprintf(LHJ_dbName, LHJ_nodeID)

	// 判断数据库是否存在
	if LHJ_DBExists(LHJ_dbName) {
		fmt.Println("创世区块已经存在!!!...")
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块======>>>>")

	// 创建或者打开数据库
	db, err := bolt.Open(LHJ_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var LHJ_genesisHash []byte

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
			LHJ_txCoinbase := LHJ_NewCoinbaseTransaction(LHJ_address)

			LHJ_genesisBlock := LHJ_CreateGenesisBlock([]*LHJ_Transaction{LHJ_txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(LHJ_genesisBlock.LHJ_Hash, LHJ_genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), LHJ_genesisBlock.LHJ_Hash)
			if err != nil {
				log.Panic(err)
			}

			LHJ_genesisHash = LHJ_genesisBlock.LHJ_Hash
		}

		return nil
	})

	return &LHJ_Blockchain{LHJ_genesisHash, db}

}

// 返回Blockchain对象
func LHJ_BlockchainObject(LHJ_nodeID string) *LHJ_Blockchain {

	LHJ_dbName := fmt.Sprintf(LHJ_dbName,LHJ_nodeID)

	// 判断数据库是否存在
	if LHJ_DBExists(LHJ_dbName) == false {
		fmt.Println("数据库不存在~~~~~")
		os.Exit(1)
	}

	db, err := bolt.Open(LHJ_dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var LHJ_tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			LHJ_tip = b.Get([]byte("l"))

		}

		return nil
	})

	return &LHJ_Blockchain{LHJ_tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (LHJ_blockchain *LHJ_Blockchain) LHJ_UnUTXOs(LHJ_address string,txs []*LHJ_Transaction) []*LHJ_UTXO {

	var LHJ_unUTXOs []*LHJ_UTXO

	LHJ_spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.LHJ_IsCoinbaseTransaction() == false {
			for _, in := range tx.LHJ_Vins {
				//是否能够解锁
				LHJ_publicKeyHash := Base58Decode([]byte(LHJ_address))

				LHJ_ripemd160Hash := LHJ_publicKeyHash[1:len(LHJ_publicKeyHash) - 4]
				if in.LHJ_UnLockRipemd160Hash(LHJ_ripemd160Hash) {

					key := hex.EncodeToString(in.LHJ_TxHash)

					LHJ_spentTXOutputs[key] = append(LHJ_spentTXOutputs[key], in.LHJ_Vout)
				}

			}
		}
	}


	for _,tx := range txs {

	Work1:
		for index,out := range tx.LHJ_Vouts {

			if out.LHJ_UnLockScriptPubKeyWithAddress(LHJ_address) {

				fmt.Println()
				fmt.Println("BY 罗浩杰...")
				fmt.Println(LHJ_address)
				fmt.Println(LHJ_spentTXOutputs)

				if len(LHJ_spentTXOutputs) == 0 {
					utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
					LHJ_unUTXOs = append(LHJ_unUTXOs, utxo)
				} else {
					for hash,indexArray := range LHJ_spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.LHJ_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
									LHJ_unUTXOs = append(LHJ_unUTXOs, utxo)
								}
							}
						} else {
							utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
							LHJ_unUTXOs = append(LHJ_unUTXOs, utxo)
						}
					}
				}
			}
		}
	}


	LHJ_blockIterator := LHJ_blockchain.LHJ_Iterator()

	for {

		block := LHJ_blockIterator.LHJ_Next()

		fmt.Println(block)
		fmt.Println()

		for i := len(block.LHJ_Txs) - 1; i >= 0 ; i-- {

			tx := block.LHJ_Txs[i]
			// txHash
			// Vins
			if tx.LHJ_IsCoinbaseTransaction() == false {
				for _, in := range tx.LHJ_Vins {
					//是否能够解锁
					LHJ_publicKeyHash := Base58Decode([]byte(LHJ_address))

					LHJ_ripemd160Hash := LHJ_publicKeyHash[1:len(LHJ_publicKeyHash) - 4]

					if in.LHJ_UnLockRipemd160Hash(LHJ_ripemd160Hash) {

						key := hex.EncodeToString(in.LHJ_TxHash)

						LHJ_spentTXOutputs[key] = append(LHJ_spentTXOutputs[key], in.LHJ_Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.LHJ_Vouts {

				if out.LHJ_UnLockScriptPubKeyWithAddress(LHJ_address) {

					fmt.Println(out)
					fmt.Println(LHJ_spentTXOutputs)

					if LHJ_spentTXOutputs != nil {

						if len(LHJ_spentTXOutputs) != 0 {

							var LHJ_isSpentUTXO bool

							for txHash, indexArray := range LHJ_spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.LHJ_TxHash) {
										LHJ_isSpentUTXO = true
										continue work
									}
								}
							}

							if LHJ_isSpentUTXO == false {

								utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
								LHJ_unUTXOs = append(LHJ_unUTXOs, utxo)

							}
						} else {
							utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
							LHJ_unUTXOs = append(LHJ_unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(LHJ_spentTXOutputs)
		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return LHJ_unUTXOs
}

// 转账时查找可用的UTXO
func (LHJ_blockchain *LHJ_Blockchain) LHJ_FindSpendableUTXOS(LHJ_from string, LHJ_amount int,LHJ_txs []*LHJ_Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	LHJ_utxos := LHJ_blockchain.LHJ_UnUTXOs(LHJ_from,LHJ_txs)

	LHJ_spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range LHJ_utxos {

		value = value + utxo.LHJ_Output.LHJ_Value

		hash := hex.EncodeToString(utxo.LHJ_TxHash)
		LHJ_spendableUTXO[hash] = append(LHJ_spendableUTXO[hash], utxo.LHJ_Index)

		if value >= int64(LHJ_amount) {
			break
		}
	}

	if value < int64(LHJ_amount) {

		fmt.Printf("%s's fund is 不足\n", LHJ_from)
		os.Exit(1)
	}

	return value, LHJ_spendableUTXO
}

// 挖掘新的区块
func (LHJ_blockchain *LHJ_Blockchain) LHJ_MineNewBlock(LHJ_from []string, LHJ_to []string, LHJ_amount []string,nodeID string) {

	//1.建立一笔交易
	LHJ_utxoSet := &LHJ_UTXOSet{LHJ_blockchain}

	var txs []*LHJ_Transaction

	for index,address := range LHJ_from {
		value, _ := strconv.Atoi(LHJ_amount[index])
		tx := LHJ_NewSimpleTransaction(address, LHJ_to[index], int64(value), LHJ_utxoSet,txs,nodeID)
		txs = append(txs, tx)
	}

	//奖励
	tx := LHJ_NewCoinbaseTransaction(LHJ_from[0])
	txs = append(txs,tx)

	//1. 通过相关算法建立Transaction数组
	var LHJ_block *LHJ_Block

	LHJ_blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			LHJ_block = DeserializeBlock(blockBytes)

		}

		return nil
	})


	// 在建立新区块之前对txs进行签名验证

	_txs := []*LHJ_Transaction{}

	for _,tx := range txs  {

		if LHJ_blockchain.LHJ_VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}


	//2. 建立新的区块
	LHJ_block = LHJ_NewBlock(txs, LHJ_block.LHJ_Height+1, LHJ_block.LHJ_Hash)

	//将新区块存储到数据库
	LHJ_blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LHJ_blockTableName))
		if b != nil {

			b.Put(LHJ_block.LHJ_Hash, LHJ_block.Serialize())

			b.Put([]byte("l"), LHJ_block.LHJ_Hash)

			LHJ_blockchain.LHJ_Tip = LHJ_block.LHJ_Hash

		}
		return nil
	})

}

// 查询余额
func (LHJ_blockchain *LHJ_Blockchain) LHJ_GetBalance(LHJ_address string) int64 {

	LHJ_utxos := LHJ_blockchain.LHJ_UnUTXOs(LHJ_address,[]*LHJ_Transaction{})

	var LHJ_amount int64

	for _, utxo := range LHJ_utxos {

		LHJ_amount = LHJ_amount + utxo.LHJ_Output.LHJ_Value
	}

	return LHJ_amount
}


func (LHJ_bclockchain *LHJ_Blockchain) LHJ_SignTransaction(LHJ_tx *LHJ_Transaction,LHJ_privKey ecdsa.PrivateKey,LHJ_txs []*LHJ_Transaction)  {

	if LHJ_tx.LHJ_IsCoinbaseTransaction() {
		return
	}

	LHJ_prevTXs := make(map[string]LHJ_Transaction)

	for _, vin := range LHJ_tx.LHJ_Vins {
		prevTX, err := LHJ_bclockchain.LHJ_FindTransaction(vin.LHJ_TxHash,LHJ_txs)
		if err != nil {
			log.Panic(err)
		}
		LHJ_prevTXs[hex.EncodeToString(prevTX.LHJ_TxHash)] = prevTX
	}

	LHJ_tx.LHJ_Sign(LHJ_privKey, LHJ_prevTXs)

}

func (LHJ_bc *LHJ_Blockchain) LHJ_FindTransaction(ID []byte,txs []*LHJ_Transaction) (LHJ_Transaction, error) {


	for _,tx := range txs  {
		if bytes.Compare(tx.LHJ_TxHash, ID) == 0 {
			return *tx, nil
		}
	}


	LHJ_bci := LHJ_bc.LHJ_Iterator()

	for {
		block := LHJ_bci.LHJ_Next()

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
func (LHJ_bc *LHJ_Blockchain) LHJ_VerifyTransaction(LHJ_tx *LHJ_Transaction,LHJ_txs []*LHJ_Transaction) bool {


	LHJ_prevTXs := make(map[string]LHJ_Transaction)

	for _, vin := range LHJ_tx.LHJ_Vins {
		prevTX, err := LHJ_bc.LHJ_FindTransaction(vin.LHJ_TxHash,LHJ_txs)
		if err != nil {
			log.Panic(err)
		}
		LHJ_prevTXs[hex.EncodeToString(prevTX.LHJ_TxHash)] = prevTX
	}

	return LHJ_tx.LHJ_Verify(LHJ_prevTXs)
}


// [string]*TXOutputs
func (LHJ_blc *LHJ_Blockchain) LHJ_FindUTXOMap() map[string]*LHJ_TXOutputs  {

	LHJ_blcIterator := LHJ_blc.LHJ_Iterator()

	// 存储已花费的UTXO的信息
	LHJ_spentableUTXOsMap := make(map[string][]*LHJ_TXInput)


	LHJ_utxoMaps := make(map[string]*LHJ_TXOutputs)


	for {
		block := LHJ_blcIterator.LHJ_Next()

		for i := len(block.LHJ_Txs) - 1; i >= 0 ;i-- {

			txOutputs := &LHJ_TXOutputs{[]*LHJ_UTXO{}}

			tx := block.LHJ_Txs[i]

			// coinbase
			if tx.LHJ_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.LHJ_Vins {

					txHash := hex.EncodeToString(txInput.LHJ_TxHash)
					LHJ_spentableUTXOsMap[txHash] = append(LHJ_spentableUTXOsMap[txHash],txInput)

				}
			}

			txHash := hex.EncodeToString(tx.LHJ_TxHash)

			txInputs := LHJ_spentableUTXOsMap[txHash]

			if len(txInputs) > 0 {


			WorkOutLoop:
				for index,out := range tx.LHJ_Vouts  {

					for _,in := range  txInputs {

						outPublicKey := out.LHJ_Ripemd160Hash
						inPublicKey := in.LHJ_PublicKey


						if bytes.Compare(outPublicKey,LHJ_Ripemd160Hash(inPublicKey)) == 0 {
							if index == in.LHJ_Vout {

								continue WorkOutLoop
							} else {

								utxo := &LHJ_UTXO{tx.LHJ_TxHash,index,out}
								txOutputs.LHJ_UTXOS = append(txOutputs.LHJ_UTXOS,utxo)
							}
						}
					}
				}

			} else {

				for index,out := range tx.LHJ_Vouts {
					utxo := &LHJ_UTXO{tx.LHJ_TxHash,index,out}
					txOutputs.LHJ_UTXOS = append(txOutputs.LHJ_UTXOS,utxo)
				}
			}


			// 设置键值对
			LHJ_utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}



	}

	return LHJ_utxoMaps
}

func (LHJ_bc *LHJ_Blockchain) LHJ_GetBestHeight() int64 {

	LHJ_block := LHJ_bc.LHJ_Iterator().LHJ_Next()

	return LHJ_block.LHJ_Height
}

func (LHJ_bc *LHJ_Blockchain) LHJ_GetBlockHashes() [][]byte {

	LHJ_blockIterator := LHJ_bc.LHJ_Iterator()

	var LHJ_blockHashs [][]byte

	for {
		block := LHJ_blockIterator.LHJ_Next()

		LHJ_blockHashs = append(LHJ_blockHashs,block.LHJ_Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.LHJ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return LHJ_blockHashs
}

func (LHJ_bc *LHJ_Blockchain) LHJ_GetBlock(LHJ_blockHash []byte) ([]byte ,error) {

	var LHJ_blockBytes []byte

	err := LHJ_bc.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {

			LHJ_blockBytes = b.Get(LHJ_blockHash)

		}

		return nil
	})

	return LHJ_blockBytes,err
}

func (LHJ_bc *LHJ_Blockchain) LHJ_AddBlock(LHJ_block *LHJ_Block)  {

	err := LHJ_bc.LHJ_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_blockTableName))

		if b != nil {

			blockExist := b.Get(LHJ_block.LHJ_Hash)

			if blockExist != nil {
				// 如果存在，不需要做任何过多的处理
				return nil
			}

			err := b.Put(LHJ_block.LHJ_Hash,LHJ_block.Serialize())

			if err != nil {
				log.Panic(err)
			}

			// 最新的区块链的Hash
			blockHash := b.Get([]byte("l"))

			blockBytes := b.Get(blockHash)

			blockInDB := DeserializeBlock(blockBytes)

			if blockInDB.LHJ_Height < LHJ_block.LHJ_Height {

				b.Put([]byte("l"),LHJ_block.LHJ_Hash)
				LHJ_bc.LHJ_Tip = LHJ_block.LHJ_Hash
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}