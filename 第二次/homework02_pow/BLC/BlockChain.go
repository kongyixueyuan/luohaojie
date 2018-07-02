package BLC

type BlockChain struct {
	Blocks []*Block // Block切片，有序存储block
}

func (bc *BlockChain) AddBlockToChain(block *Block) {
	bc.Blocks = append(bc.Blocks, block);
	//fmt.Println(len(bc.Blocks))
}
