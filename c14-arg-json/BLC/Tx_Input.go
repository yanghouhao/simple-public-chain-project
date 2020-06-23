package blc

//TxInput 交易的输出管理
type TxInput struct {
	TxHash    []byte //交易哈希（不是当前的交易哈希）
	Vout      int    //引用的上一笔交易的输出索引号
	ScriptSig string //用户名
}
