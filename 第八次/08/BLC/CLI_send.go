package BLC

import (
	"fmt"
	"strconv"
)

// 转账
func (LHJ_cli *LHJ_CLI) LHJ_send(LHJ_from []string, LHJ_to []string, LHJ_amount []string, LHJ_nodeID string, LHJ_mineNow bool) {

	blockchain := LHJ_BlockchainObject(LHJ_nodeID)
	utxoSet := &LHJ_UTXOSet{blockchain}

	defer blockchain.LHJ_DB.Close()

	if LHJ_mineNow {
		blockchain.LHJ_MineNewBlock(LHJ_from, LHJ_to, LHJ_amount, LHJ_nodeID)

		//utxoSet := &LHJ_UTXOSet{blockchain}

		//转账成功以后，需要更新一下
		utxoSet.LHJ_Update()

	} else {

		fmt.Println()
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理====>>>")

		value, _ := strconv.Atoi(LHJ_amount[0])
		tx := LHJ_NewSimpleTransaction(LHJ_from[0], LHJ_to[0], int64(value), utxoSet, []*LHJ_Transaction{}, LHJ_nodeID)
		LHJ_sendTx(LHJ_knowNodes[0], tx)

	}
}
