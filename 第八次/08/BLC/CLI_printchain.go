package BLC


func (LHJ_cli *LHJ_CLI) LHJ_printchain(LHJ_nodeID string)  {

	blockchain := LHJ_BlockchainObject(LHJ_nodeID)

	defer blockchain.LHJ_DB.Close()

	blockchain.LHJ_Printchain()

}