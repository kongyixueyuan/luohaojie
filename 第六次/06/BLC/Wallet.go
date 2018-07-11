package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"bytes"
)

const version = byte(0x00)
const addressChecksumLen = 4



type LHJ_Wallet struct {
	//1. 私钥
	LHJ_PrivateKey ecdsa.PrivateKey

	//2. 公钥
	LHJ_PublicKey  []byte
}

func IsValidForAdress(adress []byte) bool {

	// 25
	version_public_checksumBytes := Base58Decode(adress)

	fmt.Println(version_public_checksumBytes)

	//25
	//4
	//21
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes) - addressChecksumLen]

	//fmt.Println(len(checkSumBytes))
	//fmt.Println(len(version_ripemd160))

	checkBytes := LHJ_CheckSum(version_ripemd160)

	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}

	return false
}


func (w *LHJ_Wallet) LHJ_GetAddress() []byte  {

	//1. hash160
	// 20字节
	ripemd160Hash := LHJ_Ripemd160Hash(w.LHJ_PublicKey)

	// 21字节
	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)

	// 两次的256 hash
	checkSumBytes := LHJ_CheckSum(version_ripemd160Hash)

	 //25
	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Base58Encode(bytes)
}

func LHJ_CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressChecksumLen]
}


func LHJ_Ripemd160Hash(publicKey []byte) []byte {

	//1. 256

	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2. 160

	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}


// 创建钱包
func NewWallet() *LHJ_Wallet {

	privateKey,publicKey := LHJ_newKeyPair()

	return &LHJ_Wallet{privateKey,publicKey}
}


// 通过私钥产生公钥
func LHJ_newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}