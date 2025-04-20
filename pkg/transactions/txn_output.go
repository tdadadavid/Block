package transactions

// TxnOutput Outputs, transaction output, or TxOut is an output in a transaction
// which contains two fields: a value field for transferring zero or more
// satoshis and a pubkey script for indicating what conditions must be fulfilled
// for those satoshis to be further spent.
// it represents DEBIT
// Ref: https://cypherpunks-core.github.io/bitcoinbook/glossary.html
type TxnOutput struct {
	Value        int64  `json:"value"`
	ScriptPubKey string `json:"pub_key"`
}

type TxnOutputs struct {
	Outputs []TxnOutput `json:"outputs"`
}

func (to *TxnOutput) CanUnlockWith(data string) bool {
	return to.ScriptPubKey == data
}
