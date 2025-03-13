package transactions

// TxnInput Represents a transaction input
type TxnInput struct {
	txnId string
	output int32
	scriptSignature string
}

func (to *TxnInput) canUnlockWith(data string) bool {
	return to.scriptSignature == data
}