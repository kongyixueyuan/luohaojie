package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type LHJ_TXOutputs struct {
	LHJ_UTXOS []*UTXO
}


// 将区块序列化成字节数组
func (txOutputs *LHJ_TXOutputs) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeTXOutputs(txOutputsBytes []byte) *LHJ_TXOutputs {

	var txOutputs LHJ_TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}