package BLC


func (cli *CLI) resetUTXOSet(nodeID string)  {

	blockchain := BlockchainObject(nodeID)

	defer blockchain.LHJ_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.LHJ_ResetUTXOSet()

}
