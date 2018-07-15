package BLC

const LHJ_PROTOCOL  = "tcp"
const LHJ_COMMANDLENGTH  = 12
const LHJ_NODE_VERSION  = 1

// 命令
const LHJ_COMMAND_VERSION  = "version"
const LHJ_COMMAND_ADDR  = "addr"
const LHJ_COMMAND_BLOCK  = "block"
const LHJ_COMMAND_INV  = "inv"
const LHJ_COMMAND_GETBLOCKS  = "getblocks"
const LHJ_COMMAND_GETDATA  = "getdata"
const LHJ_COMMAND_TX  = "tx"

// 类型
const LHJ_BLOCK_TYPE  = "block"
const LHJ_TX_TYPE  = "tx"