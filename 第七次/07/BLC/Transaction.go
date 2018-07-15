package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"

	"math/big"
	"crypto/elliptic"
	"time"
)

// UTXO
type LHJ_Transaction struct {

	//1. 交易hash
	LHJ_TxHash []byte

	//2. 输入
	LHJ_Vins []*LHJ_TXInput

	//3. 输出
	LHJ_Vouts []*LHJ_TXOutput
}

//[]byte{}

// 判断当前的交易是否是Coinbase交易
func (tx *LHJ_Transaction) LHJ_IsCoinbaseTransaction() bool {

	return len(tx.LHJ_Vins[0].LHJ_TxHash) == 0 && tx.LHJ_Vins[0].LHJ_Vout == -1
}



//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func LHJ_NewCoinbaseTransaction(address string) *LHJ_Transaction {

	//代表消费
	txInput := &LHJ_TXInput{[]byte{},-1,nil,[]byte{}}


	txOutput := LHJ_NewTXOutput(10,address)

	txCoinbase := &LHJ_Transaction{[]byte{},[]*LHJ_TXInput{txInput},[]*LHJ_TXOutput{txOutput}}

	//设置hash值
	txCoinbase.LHJ_HashTransaction()


	return txCoinbase
}

func (tx *LHJ_Transaction) LHJ_HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	resultBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.LHJ_TxHash = hash[:]
}



//2. 转账时产生的Transaction

func NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*LHJ_Transaction,nodeID string) *LHJ_Transaction {

	//$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	wallets,_ := NewWallets(nodeID)
	wallet := wallets.LHJ_WalletsMap[from]


	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.LHJ_FindSpendableUTXOS(from,amount,txs)
	//
	//	{hash1:[0],hash2:[2,3]}

	var txIntputs []*LHJ_TXInput
	var txOutputs []*LHJ_TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &LHJ_TXInput{txHashBytes,index,nil,wallet.LHJ_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := LHJ_NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = LHJ_NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &LHJ_Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.LHJ_HashTransaction()

	//进行签名
	utxoSet.Blockchain.SignTransaction(tx, wallet.LHJ_PrivateKey,txs)

	return tx

}

func (tx *LHJ_Transaction) Hash() []byte {

	txCopy := tx

	txCopy.LHJ_TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}


func (tx *LHJ_Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}


func (tx *LHJ_Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]LHJ_Transaction) {

	if tx.LHJ_IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.LHJ_Vins {
		if prevTXs[hex.EncodeToString(vin.LHJ_TxHash)].LHJ_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.LHJ_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.LHJ_TxHash)]
		txCopy.LHJ_Vins[inID].LHJ_Signature = nil
		txCopy.LHJ_Vins[inID].LHJ_PublicKey = prevTx.LHJ_Vouts[vin.LHJ_Vout].LHJ_Ripemd160Hash
		txCopy.LHJ_TxHash = txCopy.Hash()
		txCopy.LHJ_Vins[inID].LHJ_PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.LHJ_TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.LHJ_Vins[inID].LHJ_Signature = signature
	}
}


// 拷贝一份新的Transaction用于签名                                    T
func (tx *LHJ_Transaction) TrimmedCopy() LHJ_Transaction {
	var inputs []*LHJ_TXInput
	var outputs []*LHJ_TXOutput

	for _, vin := range tx.LHJ_Vins {
		inputs = append(inputs, &LHJ_TXInput{vin.LHJ_TxHash, vin.LHJ_Vout, nil, nil})
	}

	for _, vout := range tx.LHJ_Vouts {
		outputs = append(outputs, &LHJ_TXOutput{vout.LHJ_Value, vout.LHJ_Ripemd160Hash})
	}

	txCopy := LHJ_Transaction{tx.LHJ_TxHash, inputs, outputs}

	return txCopy
}


// 数字签名验证

func (tx *LHJ_Transaction) LHJ_Verify(prevTXs map[string]LHJ_Transaction) bool {
	if tx.LHJ_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.LHJ_Vins {
		if prevTXs[hex.EncodeToString(vin.LHJ_TxHash)].LHJ_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.LHJ_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.LHJ_TxHash)]
		txCopy.LHJ_Vins[inID].LHJ_Signature = nil
		txCopy.LHJ_Vins[inID].LHJ_PublicKey = prevTx.LHJ_Vouts[vin.LHJ_Vout].LHJ_Ripemd160Hash
		txCopy.LHJ_TxHash = txCopy.Hash()
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
