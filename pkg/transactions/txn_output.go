package transactions

// TxnOutput represents a transaction output
type TxnOutput struct {
	value int
	scriptPubKey string
}

func (to *TxnOutput) canUnlockWith(data string) bool {
	return to.scriptPubKey == data
}