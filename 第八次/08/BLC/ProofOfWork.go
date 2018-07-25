package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const LHJ_targetBit  = 20
//const LHJ_targetBit  = 16


type LHJ_ProofOfWork struct {
	LHJ_Block *LHJ_Block // 当前要验证的区块
	LHJ_target *big.Int // 大数据存储
}

// 数据拼接，返回字节数组
func (LHJ_pow *LHJ_ProofOfWork) LHJ_prepareData(LHJ_nonce int, LHJ_txHash []byte) []byte {
	data := bytes.Join(
		[][]byte{
			LHJ_pow.LHJ_Block.LHJ_PrevBlockHash,
			LHJ_txHash,
			IntToHex(LHJ_pow.LHJ_Block.LHJ_Timestamp),
			IntToHex(int64(LHJ_targetBit)),
			IntToHex(int64(LHJ_nonce)),
			IntToHex(int64(LHJ_pow.LHJ_Block.LHJ_Height)),
		},
		[]byte{},
	)

	return data
}


func (LHJ_proofOfWork *LHJ_ProofOfWork) LHJ_Run() ([]byte,int64) {


	//1. 将Block的属性拼接成字节数组

	//2. 生成hash

	//3. 判断hash有效性，如果满足条件，跳出循环

	nonce := 0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	txHash := LHJ_proofOfWork.LHJ_Block.LHJ_HashTransactions()

	for {
		//准备数据
		dataBytes := LHJ_proofOfWork.LHJ_prepareData(nonce,txHash)

		// 生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x",hash)


		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		if LHJ_proofOfWork.LHJ_target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:],int64(nonce)
}


// 创建新的工作量证明对象
func LHJ_NewProofOfWork(LHJ_block *LHJ_Block) *LHJ_ProofOfWork  {

	//1. 创建一个初始值为1的target

	target := big.NewInt(1)

	//2. 左移256 - targetBit

	target = target.Lsh(target,256 - LHJ_targetBit)

	return &LHJ_ProofOfWork{LHJ_block,target}
}






