package BLC


// 创建创世区块
func (cli *CLI) createGenesisBlockchain(address string)  {

	blockchain := LHJ_CreateBlockchainWithGenesisBlock(address)
	defer blockchain.LHJ_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ResetUTXOSet()
}