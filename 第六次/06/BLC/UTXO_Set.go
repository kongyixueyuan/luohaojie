package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)



const LHJ_utxoTableName  = "LHJ_utxoTableName"

type UTXOSet struct {
	Blockchain *LHJ_Blockchain
}

// 重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet()  {

	err := utxoSet.Blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		if b != nil {


			err := tx.DeleteBucket([]byte(LHJ_utxoTableName))

			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(LHJ_utxoTableName))
		if b != nil {

			//[string]*TXOutputs
			txOutputsMap := utxoSet.Blockchain.LHJ_FindUTXOMap()


			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.Serialize())

			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (utxoSet *UTXOSet) findUTXOForAddress(address string) []*UTXO{


	var utxos []*UTXO

	utxoSet.Blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := DeserializeTXOutputs(v)

			for _,utxo := range txOutputs.UTXOS  {

				if utxo.Output.UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos,utxo)
				}
			}
		}

		return nil
	})

	return utxos
}




func (utxoSet *UTXOSet) GetBalance(address string) int64 {

	UTXOS := utxoSet.findUTXOForAddress(address)

	var amount int64

	for _,utxo := range UTXOS  {
		amount += utxo.Output.Value
	}

	return amount
}


// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *UTXOSet) FindUnPackageSpendableUTXOS(from string, txs []*LHJ_Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.LHJ_Vins {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(from))

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

			if out.UnLockScriptPubKeyWithAddress(from) {

				fmt.Println(from)

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

	return unUTXOs

}

func (utxoSet *UTXOSet) FindSpendableUTXOS(from string,amount int64,txs []*LHJ_Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {

		money += UTXO.Output.Value;
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}


	// 钱还不够
	utxoSet.Blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		if b != nil {

			c := b.Cursor()
			UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.UTXOS {

					money += utxo.Output.Value
					txHash := hex.EncodeToString(utxo.TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.Index)

					if money >= amount {
						 break UTXOBREAK;
					}
				}
			}

		}

		return nil
	})

	if money < amount{
		log.Panic("余额不足")
	}


	return  money,spentableUTXO
}


// 更新
func (utxoSet *UTXOSet) Update()  {

	// blocks
	//


	// 最新的Block
	block := utxoSet.Blockchain.Iterator().Next()


	// utxoTable
	//

	ins := []*TXInput{}

	outsMap := make(map[string]*TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.LHJ_Txs {

		for _,in := range tx.LHJ_Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.LHJ_Txs  {


		utxos := []*UTXO{}

		for index,out := range tx.LHJ_Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.Vout == index && bytes.Compare(tx.LHJ_TxHash ,in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash,LHJ_Ripemd160Hash(in.PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.LHJ_TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.LHJ_TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}

	}



	err := utxoSet.Blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		if b != nil {


			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}



				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.UTXOS  {

					if in.Vout == utxo.Index && bytes.Compare(utxo.Output.Ripemd160Hash,LHJ_Ripemd160Hash(in.PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]

						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs

					}
				}

			}

			// 新增

			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}




