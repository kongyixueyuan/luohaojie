package BLC


func (cli *CLI) printchain(nodeID string)  {

	blockchain := BlockchainObject(nodeID)

	defer blockchain.LHJ_DB.Close()

	blockchain.LHJ_Printchain()

}