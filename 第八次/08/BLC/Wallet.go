package BLC

import (
	"crypto/ecdsa"
	"fmt"
	"bytes"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

const LHJ_version = byte(0x00)
const LHJ_addressChecksumLen = 4

type LHJ_Wallet struct {
	LHJ_PrivateKey ecdsa.PrivateKey
	LHJ_PublicKey  []byte
}

func LHJ_IsValidForAddress(address []byte) bool {

	LHJ_version_public_checksumBytes := Base58Decode(address)
	fmt.Println(LHJ_version_public_checksumBytes)

	LHJ_checkSumBytes := LHJ_version_public_checksumBytes[len(LHJ_version_public_checksumBytes)-LHJ_addressChecksumLen:]

	LHJ_version_ripemd160 := LHJ_version_public_checksumBytes[:len(LHJ_version_public_checksumBytes)-LHJ_addressChecksumLen]

	LHJ_checkBytes := LHJ_CheckSum(LHJ_version_ripemd160)

	if bytes.Compare(LHJ_checkSumBytes, LHJ_checkBytes) == 0 {
		return true
	}

	return false
}

//生成一个地址
func (LHJ_w *LHJ_Wallet) LHJ_GetAddress() []byte {

	//256+160,20字节
	LHJ_ripemd160Hash := LHJ_Ripemd160Hash(LHJ_w.LHJ_PublicKey)

	//21字节
	LHJ_version_ripemd160Hash := append([]byte{LHJ_version}, LHJ_ripemd160Hash...)

	//2次256
	LHJ_checkSumBytes := LHJ_CheckSum(LHJ_version_ripemd160Hash)

	//25字节
	LHJ_bytes := append(LHJ_version_ripemd160Hash, LHJ_checkSumBytes...)

	return Base58Encode(LHJ_bytes)
}

//两次256hash后取前四位
func LHJ_CheckSum(LHJ_payload []byte) []byte {

	LHJ_hash1 := sha256.Sum256(LHJ_payload)
	LHJ_hash2 := sha256.Sum256(LHJ_hash1[:])

	return LHJ_hash2[:LHJ_addressChecksumLen]

}

//256hash->160
func LHJ_Ripemd160Hash(LHJ_publicKey []byte) []byte {
	//256
	LHJ_hash256 := sha256.New()
	LHJ_hash256.Write(LHJ_publicKey)
	LHJ_hash := LHJ_hash256.Sum(nil)

	//160
	LHJ_ripemd160 := ripemd160.New()
	LHJ_ripemd160.Write(LHJ_hash)

	return LHJ_ripemd160.Sum(nil)
}

//创建钱包
func LHJ_NewWallet() *LHJ_Wallet {
	LHJ_privateKey, LHJ_publicKey := LHJ_newKeyPair()

	return &LHJ_Wallet{LHJ_privateKey, LHJ_publicKey}
}

//私钥生成公钥
func LHJ_newKeyPair() (ecdsa.PrivateKey, []byte) {

	LHJ_curve := elliptic.P256()
	LHJ_private, err := ecdsa.GenerateKey(LHJ_curve, rand.Reader)
	if err != nil {
		log.Panic()
	}

	LHJ_pubKey := append(LHJ_private.PublicKey.X.Bytes(), LHJ_private.PublicKey.Y.Bytes()...)

	return *LHJ_private, LHJ_pubKey

}
