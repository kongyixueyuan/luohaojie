package BLC

import (
	"crypto/sha256"
)

type LHJ_MerkleTree struct {
	LHJ_RootNode *LHJ_MerkleNode
}

type LHJ_MerkleNode struct {
	LHJ_Left  *LHJ_MerkleNode
	LHJ_Right *LHJ_MerkleNode
	LHJ_Data  []byte
}

func LHJ_NewMerkleTree(LHJ_data [][]byte) *LHJ_MerkleTree {

	var LHJ_nodes []LHJ_MerkleNode

	if len(LHJ_data)%2 != 0 {
		LHJ_data = append(LHJ_data, LHJ_data[len(LHJ_data)-1])
		//[tx1,tx2,tx3,tx3]
	}

	// 创建叶子节点
	for _, datum := range LHJ_data {
		node := LHJ_NewMerkleNode(nil, nil, datum)
		LHJ_nodes = append(LHJ_nodes, *node)
	}

	// 　循环两次
	for i := 0; i < len(LHJ_data)/2; i++ {

		var newLevel []LHJ_MerkleNode

		for j := 0; j < len(LHJ_nodes); j += 2 {
			node := LHJ_NewMerkleNode(&LHJ_nodes[j], &LHJ_nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		if len(newLevel)%2 != 0 {

			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}

		LHJ_nodes = newLevel
	}

	mTree := LHJ_MerkleTree{&LHJ_nodes[0]}

	return &mTree
}

func LHJ_NewMerkleNode(LHJ_left, LHJ_right *LHJ_MerkleNode, LHJ_data []byte) *LHJ_MerkleNode {

	mNode := LHJ_MerkleNode{}

	// 创建叶子节点
	if LHJ_left == nil && LHJ_right == nil {
		hash := sha256.Sum256(LHJ_data)
		mNode.LHJ_Data = hash[:]
		// 非叶子节点
	} else {
		prevHashes := append(LHJ_left.LHJ_Data, LHJ_right.LHJ_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.LHJ_Data = hash[:]
	}

	mNode.LHJ_Left = LHJ_left
	mNode.LHJ_Right = LHJ_right

	return &mNode
}
