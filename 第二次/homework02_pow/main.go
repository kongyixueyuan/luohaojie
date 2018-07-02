package main

import (
	"homework02_pow/BLC"
	"fmt"
)

func main() {
	gBlock:=BLC.GenesisBlock("罗浩杰创世区块")
	bc := BLC.BlockChain{}
	bc.AddBlockToChain(gBlock)
	block := BLC.NewBlock(gBlock.Height+1, "交易数据一",gBlock.Hash)
	block2 := BLC.NewBlock(block.Height+1, "交易数据二",block.Hash)
	block3 := BLC.NewBlock(block2.Height+1, "交易数据三",block2.Hash)
	block4 := BLC.NewBlock(block3.Height+1, "交易数据四",block3.Hash)

	bc.AddBlockToChain(block)
	bc.AddBlockToChain(block2)
	bc.AddBlockToChain(block3)
	bc.AddBlockToChain(block4)

	fmt.Print(len(bc.Blocks))

}
