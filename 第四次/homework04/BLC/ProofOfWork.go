package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)


// 256位Hash里面前面至少要有20个零，转16进制就5个0
const targetBits = 16

/*
2018年06月28日17:11:31
工作证明算法必须满足要求：做完工作不易完成，证明工作容易完成。
 */

type ProofOfWork struct {
	Block *Block // 当前要验证的区块
	target *big.Int // 大数据存储
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)

	return data
}

/**
工作量证明，挖矿
 */
func (proofOfWork *ProofOfWork) Run() ([]byte,int64) {

	nonce := 0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)

		// 生成hash
		hash = sha256.Sum256(dataBytes)
		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])
		fmt.Printf("\r%x",hash)
		//fmt.Println(nonce)
		// Cmp compares x and y and returns
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:],int64(nonce)
}


// 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork  {

	//初始值为1的target

	target := big.NewInt(1)

	/*
	在新的工作证明的函数中，我们初始化一个值为1的big.Int，
	并将其左移256个 - targetBits位。
	 */

	target = target.Lsh(target,256 - targetBits)

	return &ProofOfWork{block,target}
}






