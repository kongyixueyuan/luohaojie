package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {

	Height int64 //区块高度

	PrevBlockHash []byte  //上一个区块HASH

	Data []byte  //交易数据

	Timestamp int64  //时间戳

	Hash []byte  //本区块Hash

	Nonce int64
}


/*创建新的区块
 */
func NewBlock(data string,height int64,prevBlockHash []byte) *Block {

	//根据交易数据，区块高度，上一区块Hash创建区块数据
	block := &Block{height,prevBlockHash,[]byte(data),time.Now().Unix(),nil,0}

	// 初始化工作量证明，并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	// 挖矿验证
	hash,nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block

}

/*创世区块
 */
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}



/*将区块序列化成字节数组存储到数据库中
 */
func (block *Block) Serialize() []byte {

	//创建缓存区
	var resultBuffer bytes.Buffer
	encoder := gob.NewEncoder(&resultBuffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return resultBuffer.Bytes()
}

/*从数据库中查询出数据，将数据反序列化城对象
 */
func DeserializeBlock(blockBytes []byte) *Block {

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
