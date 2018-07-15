package BLC

import "bytes"


type LHJ_TXOutput struct {
	LHJ_Value int64
	LHJ_Ripemd160Hash []byte  //用户名
}

func (txOutput *LHJ_TXOutput)  LHJ_Lock(address string)  {

	publicKeyHash := Base58Decode([]byte(address))

	txOutput.LHJ_Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}


func LHJ_NewTXOutput(value int64,address string) *LHJ_TXOutput {

	txOutput := &LHJ_TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.LHJ_Lock(address)

	return txOutput
}


// 解锁
func (txOutput *LHJ_TXOutput) LHJ_UnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(txOutput.LHJ_Ripemd160Hash,hash160) == 0
}



