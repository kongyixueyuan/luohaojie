package BLC

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

const LHJ_walletFile = "LHJ_Wallets_%s.dat"

type LHJ_Wallets struct {
	LHJ_WalletsMap map[string]*LHJ_Wallet
}

//创建钱包集合
func LHJ_NewWallets(LHJ_nodeID string) (*LHJ_Wallets, error) {

	LHJ_walletFile := fmt.Sprintf(LHJ_walletFile, LHJ_nodeID)

	if _, err := os.Stat(LHJ_walletFile); os.IsNotExist(err) {
		LHJ_wallets := &LHJ_Wallets{}
		LHJ_wallets.LHJ_WalletsMap = make(map[string]*LHJ_Wallet)
		return LHJ_wallets, err
	}

	LHJ_fileContent, err := ioutil.ReadFile(LHJ_walletFile)
	if err != nil {
		log.Panic(err)
	}

	var LHJ_wallets LHJ_Wallets
	gob.Register(elliptic.P256())
	LHJ_decoder := gob.NewDecoder(bytes.NewReader(LHJ_fileContent))
	err = LHJ_decoder.Decode(&LHJ_wallets)
	if err != nil {
		log.Panic(err)
	}

	return &LHJ_wallets, nil

}

//创建新钱包
func (LHJ_w *LHJ_Wallets) LHJ_CreateNewWallet(LHJ_nodeID string) {
	LHJ_wallet := LHJ_NewWallet()
	fmt.Printf("钱包地址: %s\n", LHJ_wallet.LHJ_GetAddress())
	LHJ_w.LHJ_WalletsMap[string(LHJ_wallet.LHJ_GetAddress())] = LHJ_wallet
	LHJ_w.LHJ_SaveWallets(LHJ_nodeID)
}

//将钱包信息写入到文件
func (LHJ_w *LHJ_Wallets) LHJ_SaveWallets(LHJ_nodeID string) {
	LHJ_walletFile := fmt.Sprintf(LHJ_walletFile, LHJ_nodeID)

	var LHJ_content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&LHJ_content)
	err := encoder.Encode(&LHJ_w)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件,原来文件的数据会被覆盖
	err = ioutil.WriteFile(LHJ_walletFile, LHJ_content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
