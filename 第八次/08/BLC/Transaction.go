package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
	"fmt"
	"encoding/json"
)

type LHJ_Transaction struct {
	LHJ_TxHash []byte
	LHJ_Vins   []*LHJ_TXInput
	LHJ_Vouts  []*LHJ_TXOutput
}

func (LHJ_tx *LHJ_Transaction) LHJ_PrintTX() {

	fmt.Printf("txHash : %s\n", hex.EncodeToString(LHJ_tx.LHJ_TxHash))

	fmt.Println("Vins====")
	for _, vin := range LHJ_tx.LHJ_Vins {
		vin.LHJ_PrintInfo()
	}
	fmt.Println("Vouts====")
	for _, vout := range LHJ_tx.LHJ_Vouts {
		vout.LHJ_PrintInfo()
	}
	fmt.Println("============================")
}

// 判断当前的交易是否是Coinbase交易
func (tx *LHJ_Transaction) LHJ_IsCoinbaseTransaction() bool {

	return len(tx.LHJ_Vins[0].LHJ_TxHash) == 0 && tx.LHJ_Vins[0].LHJ_Vout == -1
}

//Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func LHJ_NewCoinbaseTransaction(address string) *LHJ_Transaction {

	//代表消费
	LHJ_txInput := &LHJ_TXInput{[]byte{}, -1, nil, []byte{}}

	LHJ_txOutput := LHJ_NewTXOutput(10, address)

	LHJ_txCoinbase := &LHJ_Transaction{[]byte{}, []*LHJ_TXInput{LHJ_txInput}, []*LHJ_TXOutput{LHJ_txOutput}}

	//设置hash
	LHJ_txCoinbase.LHJ_HashTransaction()

	return LHJ_txCoinbase
}

func (LHJ_tx *LHJ_Transaction) LHJ_HashTransaction() {

	var LHJ_result bytes.Buffer

	LHJ_encoder := gob.NewEncoder(&LHJ_result)

	err := LHJ_encoder.Encode(LHJ_tx)
	if err != nil {
		log.Panic(err)
	}

	resultBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), LHJ_result.Bytes()}, []byte{})

	hash := sha256.Sum256(resultBytes)

	LHJ_tx.LHJ_TxHash = hash[:]
}

//2. 转账时产生的Transaction
func LHJ_NewSimpleTransaction(LHJ_from string,LHJ_to string,LHJ_amount int64,LHJ_utxoSet *LHJ_UTXOSet,LHJ_txs []*LHJ_Transaction,LHJ_nodeID string) *LHJ_Transaction {

	wallets,_ := LHJ_NewWallets(LHJ_nodeID)
	wallet := wallets.LHJ_WalletsMap[LHJ_from]

	// 通过一个函数，返回
	LHJ_money,LHJ_spendableUTXODic := LHJ_utxoSet.LHJ_FindSpendableUTXOS(LHJ_from,LHJ_amount,LHJ_txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*LHJ_TXInput
	var txOutputs []*LHJ_TXOutput

	for txHash,indexArray := range LHJ_spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &LHJ_TXInput{txHashBytes,index,nil,wallet.LHJ_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := LHJ_NewTXOutput(int64(LHJ_amount),LHJ_to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = LHJ_NewTXOutput(int64(LHJ_money) - int64(LHJ_amount),LHJ_from)
	txOutputs = append(txOutputs,txOutput)

	fmt.Println("找零:",int64(LHJ_money) ,"-",int64(LHJ_amount),":",int64(LHJ_money) - int64(LHJ_amount))


	tx := &LHJ_Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.LHJ_HashTransaction()

	//进行签名
	LHJ_utxoSet.LHJ_Blockchain.LHJ_SignTransaction(tx, wallet.LHJ_PrivateKey,LHJ_txs)

	tx.LHJ_PrintTX()

	return tx

}

func (LHJ_tx *LHJ_Transaction) LHJ_Hash() []byte {

	LHJ_txCopy := LHJ_tx

	LHJ_txCopy.LHJ_TxHash = []byte{}

	hash := sha256.Sum256(LHJ_txCopy.LHJ_Serialize())
	return hash[:]
}


func (LHJ_tx *LHJ_Transaction) LHJ_Serialize() []byte {
	jsonByte,err := json.Marshal(LHJ_tx)
	if err != nil{
		//fmt.Println("序列化失败:",err)
		log.Panic(err)
	}
	return jsonByte
}


func (LHJ_tx *LHJ_Transaction) LHJ_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]LHJ_Transaction) {

	if LHJ_tx.LHJ_IsCoinbaseTransaction() {
		return
	}


	for _, vin := range LHJ_tx.LHJ_Vins {
		if prevTXs[hex.EncodeToString(vin.LHJ_TxHash)].LHJ_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	LHJ_txCopy := LHJ_tx.LHJ_TrimmedCopy()

	fmt.Println("签名前")

	for inID, vin := range LHJ_txCopy.LHJ_Vins {

		fmt.Println("签名中")

		prevTx := prevTXs[hex.EncodeToString(vin.LHJ_TxHash)]
		LHJ_txCopy.LHJ_Vins[inID].LHJ_Signature = nil
		LHJ_txCopy.LHJ_Vins[inID].LHJ_PublicKey = prevTx.LHJ_Vouts[vin.LHJ_Vout].LHJ_Ripemd160Hash
		LHJ_txCopy.LHJ_TxHash = LHJ_txCopy.LHJ_Hash()
		LHJ_txCopy.LHJ_Vins[inID].LHJ_PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, LHJ_txCopy.LHJ_TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		LHJ_tx.LHJ_Vins[inID].LHJ_Signature = signature
	}
}


// 拷贝一份新的Transaction用于签名                                    T
func (LHJ_tx *LHJ_Transaction) LHJ_TrimmedCopy() LHJ_Transaction {
	var LHJ_inputs []*LHJ_TXInput
	var LHJ_outputs []*LHJ_TXOutput

	for _, vin := range LHJ_tx.LHJ_Vins {
		LHJ_inputs = append(LHJ_inputs, &LHJ_TXInput{vin.LHJ_TxHash, vin.LHJ_Vout, nil, nil})
	}

	for _, vout := range LHJ_tx.LHJ_Vouts {
		LHJ_outputs = append(LHJ_outputs, &LHJ_TXOutput{vout.LHJ_Value, vout.LHJ_Ripemd160Hash})
	}

	txCopy := LHJ_Transaction{LHJ_tx.LHJ_TxHash, LHJ_inputs, LHJ_outputs}

	return txCopy
}


// 数字签名验证

func (LHJ_tx *LHJ_Transaction) LHJ_Verify(prevTXs map[string]LHJ_Transaction) bool {

	//fmt.Println("数字签名验证:TxHash")
	//
	//fmt.Println(LHJ_tx.LHJ_TxHash)
	//
	//fmt.Println("签名.")
	//
	//fmt.Println(LHJ_tx.LHJ_Vins[0].LHJ_Signature)
	//
	//fmt.Println("公钥.")
	//
	//fmt.Println(LHJ_tx.LHJ_Vouts[0].LHJ_Ripemd160Hash)

	if LHJ_tx.LHJ_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range LHJ_tx.LHJ_Vins {
		if prevTXs[hex.EncodeToString(vin.LHJ_TxHash)].LHJ_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := LHJ_tx.LHJ_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range LHJ_tx.LHJ_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.LHJ_TxHash)]
		txCopy.LHJ_Vins[inID].LHJ_Signature = nil
		txCopy.LHJ_Vins[inID].LHJ_PublicKey = prevTx.LHJ_Vouts[vin.LHJ_Vout].LHJ_Ripemd160Hash
		txCopy.LHJ_TxHash = txCopy.LHJ_Hash()
		txCopy.LHJ_Vins[inID].LHJ_PublicKey = nil


		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.LHJ_Signature)
		r.SetBytes(vin.LHJ_Signature[:(sigLen / 2)])
		s.SetBytes(vin.LHJ_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.LHJ_PublicKey)
		x.SetBytes(vin.LHJ_PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.LHJ_PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.LHJ_TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}