package BLC

type UTXO struct {
	TxHash []byte  //txHash
	Index int      //txOutput的下标
	Output *TXOutput
}





