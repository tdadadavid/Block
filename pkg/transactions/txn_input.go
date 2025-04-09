package transactions

// TxnInput Represents a transaction input
// It represents CREDIT
type TxnInput struct {
	TxnId           string `json:"id"`
	Output          int32  `json:"out"`
	ScriptSignature string `json:"signature"`
}

type TxnInputs struct {
	Inputs []TxnInput `json:"inputs"`
}

func (to *TxnInput) CanUnlockWith(data string) bool {
	return to.ScriptSignature == data
}
