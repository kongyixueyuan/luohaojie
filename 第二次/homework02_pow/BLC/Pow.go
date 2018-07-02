package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// 256位Hash里面前面至少要有20个零，转16进制就5个0
const targetBits = 20

/*
2018年06月28日17:11:31
工作证明算法必须满足要求：做完工作不易完成，证明工作容易完成。
 */
type Pow struct {
	Block  *Block
	target *big.Int
}

//创建个工作证明对象
func NewPow(block *Block) *Pow {
	//初始化 target为1
	target := big.NewInt(1)
	/*
	在新的工作证明的函数中，我们初始化一个值为1的big.Int，
	并将其左移256个 - targetBits位。
	 */
	target = target.Lsh(target, 256-targetBits)
	return &Pow{block, target}

}

// 数据拼接，返回字节数组

func (pow *Pow) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			IntToHex(int64(pow.Block.Height)),
			pow.Block.TxData,
			pow.Block.PHash,
			IntToHex(pow.Block.CreateTime),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//挖矿
func (pow *Pow) Run() ([]byte, int64) {
	nonce := 0
	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte
	for {
		//准备数据
		dataBytes := pow.prepareData(nonce)

		// 生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce++
	}
	return hash[:], int64(nonce)

}
