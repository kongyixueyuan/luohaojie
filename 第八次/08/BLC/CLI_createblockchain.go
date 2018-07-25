package BLC

//创建创世区块
func (LHJ_cli *LHJ_CLI) LHJ_createGenesisBlockchain(LHJ_address string, LHJ_nodeID string) {

	blockchain := LHJ_CreateBlockchainWithGenesisBlock(LHJ_address, LHJ_nodeID)

	defer blockchain.LHJ_DB.Close()

	utxoSet := &LHJ_UTXOSet{blockchain}

	utxoSet.LHJ_ResetUTXOSet()
}
