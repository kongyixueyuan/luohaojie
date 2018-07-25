package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)



const LHJ_utxoTableName  = "LHJ_utxoTableName"

type LHJ_UTXOSet struct {
	LHJ_Blockchain *LHJ_Blockchain
}

// 重置数据库表
func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_ResetUTXOSet()  {

	err := LHJ_utxoSet.LHJ_Blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error {

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
			txOutputsMap := LHJ_utxoSet.LHJ_Blockchain.LHJ_FindUTXOMap()

			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.Serialize())

			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("重置失败......")
		log.Panic(err)
	}

}

func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_findUTXOForAddress(LHJ_address string) []*LHJ_UTXO{


	var LHJ_utxos []*LHJ_UTXO

	LHJ_utxoSet.LHJ_Blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := DeserializeTXOutputs(v)

			for _,utxo := range txOutputs.LHJ_UTXOS  {

				if utxo.LHJ_Output.LHJ_UnLockScriptPubKeyWithAddress(LHJ_address) {
					LHJ_utxos = append(LHJ_utxos,utxo)
				}
			}
		}

		return nil
	})

	return LHJ_utxos
}




func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_GetBalance(LHJ_address string) int64 {

	LHJ_UTXOS := LHJ_utxoSet.LHJ_findUTXOForAddress(LHJ_address)

	var LHJ_amount int64

	for _,utxo := range LHJ_UTXOS  {
		LHJ_amount += utxo.LHJ_Output.LHJ_Value
	}

	return LHJ_amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_FindUnPackageSpendableUTXOS(LHJ_from string, LHJ_txs []*LHJ_Transaction) []*LHJ_UTXO {

	var LHJ_unUTXOs []*LHJ_UTXO

	LHJ_spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range LHJ_txs {

		if tx.LHJ_IsCoinbaseTransaction() == false {
			for _, in := range tx.LHJ_Vins {
				//是否能够解锁
				LHJ_publicKeyHash := Base58Decode([]byte(LHJ_from))

				LHJ_ripemd160Hash := LHJ_publicKeyHash[1:len(LHJ_publicKeyHash) - 4]
				if in.LHJ_UnLockRipemd160Hash(LHJ_ripemd160Hash) {

					key := hex.EncodeToString(in.LHJ_TxHash)

					LHJ_spentTXOutputs[key] = append(LHJ_spentTXOutputs[key], in.LHJ_Vout)
				}
			}
		}
	}

	for _,tx := range LHJ_txs {

	Work1:
		for index,out := range tx.LHJ_Vouts {

			if out.LHJ_UnLockScriptPubKeyWithAddress(LHJ_from) {


				if len(LHJ_spentTXOutputs) == 0 {
					LHJ_utxo := &LHJ_UTXO{tx.LHJ_TxHash, index, out}
					LHJ_unUTXOs = append(LHJ_unUTXOs, LHJ_utxo)
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

	return LHJ_unUTXOs

}

func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_FindSpendableUTXOS(LHJ_from string,LHJ_amount int64,LHJ_txs []*LHJ_Transaction) (int64,map[string][]int)  {

	LHJ_unPackageUTXOS := LHJ_utxoSet.LHJ_FindUnPackageSpendableUTXOS(LHJ_from,LHJ_txs)

	LHJ_spentableUTXO := make(map[string][]int)

	var LHJ_money int64 = 0

	for _,UTXO := range LHJ_unPackageUTXOS {

		LHJ_money += UTXO.LHJ_Output.LHJ_Value;
		LHJ_txHash := hex.EncodeToString(UTXO.LHJ_TxHash)
		LHJ_spentableUTXO[LHJ_txHash] = append(LHJ_spentableUTXO[LHJ_txHash],UTXO.LHJ_Index)
		if LHJ_money >= LHJ_amount{
			return  LHJ_money,LHJ_spentableUTXO
		}
	}

	// 钱还不够
	LHJ_utxoSet.LHJ_Blockchain.LHJ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		if b != nil {

			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.LHJ_UTXOS {
					//判断是否是自己的余额
					if utxo.LHJ_Output.LHJ_UnLockScriptPubKeyWithAddress(LHJ_from) {
						LHJ_money += utxo.LHJ_Output.LHJ_Value
						txHash := hex.EncodeToString(utxo.LHJ_TxHash)
						LHJ_spentableUTXO[txHash] = append(LHJ_spentableUTXO[txHash], utxo.LHJ_Index)

						if LHJ_money >= LHJ_amount {
							break UTXOBREAK;
						}
					}
				}
			}
		}

		return nil
	})

	if LHJ_money < LHJ_amount{
		log.Panic("余额不足......")
	}


	return  LHJ_money,LHJ_spentableUTXO
}


// 更新
func (LHJ_utxoSet *LHJ_UTXOSet) LHJ_Update()  {
	// 最新的Block
	block := LHJ_utxoSet.LHJ_Blockchain.LHJ_Iterator().LHJ_Next()

	// utxoTable

	ins := []*LHJ_TXInput{}

	outsMap := make(map[string]*LHJ_TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.LHJ_Txs {

		for _,in := range tx.LHJ_Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.LHJ_Txs  {

		utxos := []*LHJ_UTXO{}

		for index,out := range tx.LHJ_Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.LHJ_Vout == index && bytes.Compare(tx.LHJ_TxHash ,in.LHJ_TxHash) == 0 && bytes.Compare(out.LHJ_Ripemd160Hash,LHJ_Ripemd160Hash(in.LHJ_PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				LHJ_utxo := &LHJ_UTXO{tx.LHJ_TxHash,index,out}
				utxos = append(utxos,LHJ_utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.LHJ_TxHash)
			outsMap[txHash] = &LHJ_TXOutputs{utxos}
		}
	}



	err := LHJ_utxoSet.LHJ_Blockchain.LHJ_DB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(LHJ_utxoTableName))

		if b != nil {

			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.LHJ_TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println("DeserializeTXOutputs")
				fmt.Println(txOutputsBytes)

				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*LHJ_UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.LHJ_UTXOS  {

					if in.LHJ_Vout == utxo.LHJ_Index && bytes.Compare(utxo.LHJ_Output.LHJ_Ripemd160Hash,LHJ_Ripemd160Hash(in.LHJ_PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}


				if isNeedDelete {
					b.Delete(in.LHJ_TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.LHJ_TxHash)]

						preTXOutputs.LHJ_UTXOS = append(preTXOutputs.LHJ_UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.LHJ_TxHash)] = preTXOutputs

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