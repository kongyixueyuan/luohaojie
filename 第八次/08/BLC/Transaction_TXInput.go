package BLC

import (
	"bytes"
	"fmt"
	"encoding/hex"
)

type LHJ_TXInput struct {
	LHJ_TxHash    []byte
	LHJ_Vout      int
	LHJ_Signature []byte
	LHJ_PublicKey []byte
}

//判断当前消费是谁的钱
func (LHJ_txInput *LHJ_TXInput) LHJ_UnLockRipemd160Hash(LHJ_ripemd160Hash []byte) bool {

	LHJ_publicKey := LHJ_Ripemd160Hash(LHJ_txInput.LHJ_PublicKey)

	return bytes.Compare(LHJ_publicKey, LHJ_ripemd160Hash) == 0
}

func (LHJ_txInput *LHJ_TXInput) LHJ_PrintInfo()  {
	fmt.Printf("txHash:%s\n",hex.EncodeToString(LHJ_txInput.LHJ_TxHash))
	fmt.Printf("Vout:%d\n",LHJ_txInput.LHJ_Vout)
	fmt.Printf("Signature:%v\n",LHJ_txInput.LHJ_Signature)
	fmt.Printf("PublicKey:%v\n",LHJ_txInput.LHJ_PublicKey)
}