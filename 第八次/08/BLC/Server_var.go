package BLC


//存储节点全局变量

//localhost:3000 主节点的地址
var LHJ_knowNodes = []string{"localhost:3000"}
var LHJ_nodeAddress string //全局变量，节点地址

// 存储hash值
var LHJ_transactionArray [][]byte
var LHJ_minerAddress string
var LHJ_memoryTxPool = make(map[string]*LHJ_Transaction)
