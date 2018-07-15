package BLC

import "bytes"

type LHJ_TXInput struct {
	// 1. 交易的Hash
	LHJ_TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	LHJ_Vout      int

	LHJ_Signature []byte // 数字签名

	LHJ_PublicKey    []byte // 公钥，钱包里面
}



// 判断当前的消费是谁的钱
func (txInput *LHJ_TXInput) LHJ_UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ripemd160Hash(txInput.LHJ_PublicKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}