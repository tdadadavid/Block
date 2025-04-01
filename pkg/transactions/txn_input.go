package transactions

// TxnInput Represents a transaction input
type TxnInput struct {
	TxnId           string `json:"id"`
	Output          int32  `json:"out"`
	ScriptSignature string `json:"signature"`
}

func (to *TxnInput) CanUnlockWith(data string) bool {
	return to.ScriptSignature == data
}
