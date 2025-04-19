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

func (ti *TxnInput) CanUnlockWith(data string) bool {
	return ti.ScriptSignature == data
}
