package BLC

import (
	"crypto/sha256"
)


type MerkleTree struct {
	RootNode *MerkleNode
}


// Block  [tx1 tx2 tx3 tx3]


//MerkleNode{nil,nil,tx1Bytes}
//MerkleNode{nil,nil,tx2Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//
//

//
//MerkleNode:
//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
//
//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
//
//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))





type MerkleNode struct {
	LHJ_Left  *MerkleNode
	LHJ_Right *MerkleNode
	LHJ_Data  []byte
}


func LHJ_NewMerkleTree(data [][]byte) *MerkleTree {

	//[tx1,tx2,tx3]

	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
		//[tx1,tx2,tx3,tx3]
	}

	// 创建叶子节点
	for _, datum := range data {
		node := LHJ_NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}


	//MerkleNode{nil,nil,tx1Bytes}
	//MerkleNode{nil,nil,tx2Bytes}
	//MerkleNode{nil,nil,tx3Bytes}
	//MerkleNode{nil,nil,tx3Bytes}



	// 　循环两次
	for i := 0; i < len(data)/2; i++ {

		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := LHJ_NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		//MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
		//
		//MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
		//


		nodes = newLevel
	}

	//MerkleNode:
	//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
	//
	//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
	//
	//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}


func LHJ_NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.LHJ_Data = hash[:]
	} else {
		prevHashes := append(left.LHJ_Data, right.LHJ_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.LHJ_Data = hash[:]
	}

	mNode.LHJ_Left = left
	mNode.LHJ_Right = right

	return &mNode
}