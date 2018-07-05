package BLC



type TXOutput struct {
	Value int64		//
	ScriptPubKey string  //用户名、公钥
}

// 解锁 ScriptPubKey是否当前的查询的人
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {

	return txOutput.ScriptPubKey == address
}


