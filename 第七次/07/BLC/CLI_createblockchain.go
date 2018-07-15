package BLC


// 创建创世区块
func (cli *CLI) createGenesisBlockchain(address string,nodeID string)  {

	blockchain := LHJ_CreateBlockchainWithGenesisBlock(address,nodeID)
	defer blockchain.LHJ_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.LHJ_ResetUTXOSet()
}

//blocks
//utxoTable