package transactions

// TxnOutput represents a transaction output
type TxnOutput struct {
	Value        int64  `json:"value"`
	ScriptPubKey string `json:"pub_key"`
}

func (to *TxnOutput) CanUnlockWith(data string) bool {
	return to.ScriptPubKey == data
}
