package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"log"
	"os"
)

const walletFile  = "LHJ_Wallets.dat"

type LHJ_Wallets struct {
	LHJ_WalletsMap map[string]*LHJ_Wallet
}



// 创建钱包集合
func LHJ_NewWallets() (*LHJ_Wallets,error){

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &LHJ_Wallets{}
		wallets.LHJ_WalletsMap = make(map[string]*LHJ_Wallet)
		return wallets,err
	}


	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets LHJ_Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}





// 创建一个新钱包
func (w *LHJ_Wallets) LHJ_CreateNewWallet()  {

	wallet := NewWallet()
	fmt.Printf("Address：%s\n",wallet.LHJ_GetAddress())
	w.LHJ_WalletsMap[string(wallet.LHJ_GetAddress())] = wallet
	w.LHJ_SaveWallets()
}

// 将钱包信息写入到文件
func (w *LHJ_Wallets) LHJ_SaveWallets()  {
	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}


}

