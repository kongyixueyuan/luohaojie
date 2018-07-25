package BLC

func (LHJ_cli *LHJ_CLI) LHJ_resetUTXOSet(LHJ_nodeID string)  {

	blockchain := LHJ_BlockchainObject(LHJ_nodeID)

	defer blockchain.LHJ_DB.Close()

	utxoSet := &LHJ_UTXOSet{blockchain}

	utxoSet.LHJ_ResetUTXOSet()

}
