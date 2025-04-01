package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactions_NewCoinbase(t *testing.T) {
	coinbase := NewCoinbase("0x00", "data")
	assert.NotNil(t, coinbase)

	assert.Empty(t, coinbase.Id)
	assert.Empty(t, coinbase.Inputs[0].TxnId)
	assert.Equal(t, coinbase.Inputs[0].Output, int32(-1))
}

func TestTransactions_Serialize_Deserialize(t *testing.T) {
	txn := Transaction{
		Id: "tx123",
		Inputs: []TxnInput{
			{TxnId: "prevTxn1", Output: 0, ScriptSignature: "sig1"},
		},
		Outputs: []TxnOutput{
			{Value: 100, ScriptPubKey: "pubKey1"},
		},
	}

	bytes, err := txn.Serialize()
	assert.NoError(t, err)

	txn1 := Transaction{}
	err = txn1.Deserialize(bytes)
	assert.NoError(t, err)

	assert.Equal(t, txn.GetId(), txn1.GetId())
	assert.ElementsMatch(t, txn.GetInputs(), txn1.GetInputs())
	assert.ElementsMatch(t, txn.GetOutputs(), txn1.GetOutputs())
}

func TestTransactions_IsCoinBase(t *testing.T) {
	tests := map[string]struct {
		txn  Transaction
		itIs bool
	}{
		"coinbase transactions": {
			txn: Transaction{
				Id: "",
				Inputs: []TxnInput{
					{
						TxnId:           "",
						Output:          -1,
						ScriptSignature: "",
					},
				},
				Outputs: []TxnOutput{
					{
						Value:        100,
						ScriptPubKey: "",
					},
				},
			},
			itIs: true,
		},
		"not coinbase transactions": {
			txn:  Transaction{},
			itIs: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			result := tc.txn.IsCoinbase()
			if result != tc.itIs {
				t.Errorf("error: expected(%v) given(%v)", result, tc.itIs)
			}
		})
	}
}
