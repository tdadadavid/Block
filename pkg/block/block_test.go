package block

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tdadadavid/block/pkg/transactions"
)

func TestBlock_New(t *testing.T) {
	b := New(transactions.Transaction{}, "", 4)

	assert.NotNil(t, b)
	assert.Equal(t, b.GetHeight(), int32(4))
	assert.NotNil(t, b.GetTimestamp())
	assert.NotNil(t, b.GetHash())
	assert.Equal(t, b.GetTimestamp(), time.Now().Unix())
	assert.Empty(t, b.GetPrevBlockHash())
	assert.NotNil(t, b.GetPrevBlockHash())
	assert.NotNil(t, b.GetTransaction())
	assert.Len(t, b.GetTransaction(), 1)
}

func TestBlock_Serialize_Deserialize(t *testing.T) {

	txn := transactions.Transaction{
		Id: "tx123",
		Inputs: []transactions.TxnInput{
			{TxnId: "prevTxn1", Output: 0, ScriptSignature: "sig1"},
		},
		Outputs: []transactions.TxnOutput{
			{Value: 100, ScriptPubKey: "pubKey1"},
		},
	}

	b := New(txn, "", 4)

	bytes, err := b.Serialize()
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	b2 := Block{}
	err = b2.Deserialize(bytes)
	assert.NoError(t, err)
	// we cannot compare equality by using assert.Equal() because of the logger in the Block struct

	assert.Equal(t, b.GetHash(), b2.GetHash())
	assert.Equal(t, b.GetHeight(), b2.GetHeight())
	assert.Equal(t, b.GetNonce(), b2.GetNonce())
	assert.Equal(t, b.GetPrevBlockHash(), b2.GetPrevBlockHash())
	assert.Equal(t, b.GetTransaction(), b2.GetTransaction())
	assert.Equal(t, b.GetTimestamp(), b2.GetTimestamp())
}

func TestBlock_NewGenesisBlock(t *testing.T) {
	txn := transactions.Transaction{
		Id:      "",
		Inputs:  []transactions.TxnInput{},
		Outputs: []transactions.TxnOutput{},
	}
	g := NewGenesisBlock(txn)

	assert.NotNil(t, g)
	assert.Equal(t, g.GetHeight(), int32(0))
	assert.NotNil(t, g.GetTimestamp())
	assert.Empty(t, g.GetPrevBlockHash())
	assert.NotNil(t, g.GetTransaction())
	assert.Equal(t, len(g.GetTransaction()), 1)
}
