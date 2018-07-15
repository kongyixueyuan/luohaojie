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

const LHJ_walletFile  = "Wallets_%s.dat"

type LHJ_Wallets struct {
	LHJ_WalletsMap map[string]*LHJ_Wallet
}



// 创建钱包集合
func NewWallets(nodeID string) (*LHJ_Wallets,error){

	walletFile := fmt.Sprintf(LHJ_walletFile,nodeID)

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
func (w *LHJ_Wallets) CreateNewWallet(nodeID string)  {

	wallet := NewWallet()
	fmt.Printf("Address：%s\n",wallet.GetAddress())
	w.LHJ_WalletsMap[string(wallet.GetAddress())] = wallet
	w.SaveWallets(nodeID)
}

// 将钱包信息写入到文件
func (w *LHJ_Wallets) SaveWallets(nodeID string)  {


	walletFile := fmt.Sprintf(LHJ_walletFile,nodeID)

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

