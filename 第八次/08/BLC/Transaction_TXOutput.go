package BLC

import (
	"bytes"
	"fmt"
	"encoding/hex"
)

type LHJ_TXOutput struct {
	LHJ_Value         int64
	LHJ_Ripemd160Hash []byte
}

func (LHJ_txOutput *LHJ_TXOutput) LHJ_Lock(LHJ_address string) {
	LHJ_publicKeyHash := Base58Decode([]byte(LHJ_address))
	LHJ_txOutput.LHJ_Ripemd160Hash = LHJ_publicKeyHash[1 : len(LHJ_publicKeyHash)-4]
}

func LHJ_NewTXOutput(LHJ_value int64, LHJ_address string) *LHJ_TXOutput {

	LHJ_txOutput := &LHJ_TXOutput{LHJ_value, nil}

	//设置Ripemd160Hash
	LHJ_txOutput.LHJ_Lock(LHJ_address)

	return LHJ_txOutput
}

//解锁
func (LHJ_txOutput *LHJ_TXOutput) LHJ_UnLockScriptPubKeyWithAddress(LHJ_address string) bool {

	LHJ_publicKeyHash := Base58Decode([]byte(LHJ_address))

	LHJ_hash160 := LHJ_publicKeyHash[1 : len(LHJ_publicKeyHash)-4]

	return bytes.Compare(LHJ_txOutput.LHJ_Ripemd160Hash,LHJ_hash160) == 0
}

func (LHJ_txOutput *LHJ_TXOutput) LHJ_PrintInfo()  {
	fmt.Printf("Value:%s\n",LHJ_txOutput.LHJ_Value)
	fmt.Printf("Ripemd160Hash:%s\n",hex.EncodeToString(LHJ_txOutput.LHJ_Ripemd160Hash))
}