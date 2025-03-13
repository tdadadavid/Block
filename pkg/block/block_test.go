package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_New(t *testing.T) {
	b := New("data", "", 4)
	
	assert.NotNil(t, b)
	assert.Equal(t, b.Height, int32(4))
	assert.NotNil(t, b.Timestamp)
	assert.Empty(t, b.PrevBlockHash)

	assert.Equal(t, b.GetHash(), b.Hash)
	assert.Equal(t, b.GetTransaction(), b.Transactions)
	assert.Equal(t, b.GetPrevBlockHash(), b.PrevBlockHash)
}

func TestBlock_Serialize_Deserialize(t *testing.T) {
	b := New("data", "", 4)

	bytez, err := b.Serialize()
	assert.Nil(t, err)
	assert.NotNil(t, bytez)

	var b2 Block
	b2.Deserialize(bytez)

	assert.Equal(t, b.Hash, b2.Hash)
	assert.Equal(t, b.Height, b2.Height)
	assert.Equal(t, b.Nonce, b2.Nonce)
	assert.Equal(t, b.PrevBlockHash, b2.PrevBlockHash)
	assert.Equal(t, b.Transactions, b2.Transactions)
	assert.Equal(t, b.Timestamp, b2.Timestamp)
}

func TestBlock_NewGenesisBlock(t *testing.T) {
	g := NewGenesisBlock()

	assert.NotNil(t, g)
	assert.Equal(t, g.Height, int32(0))
	assert.Equal(t, g.Nonce, int32(0))
	assert.NotNil(t, g.Timestamp)
	assert.Empty(t, g.PrevBlockHash)
	assert.Equal(t, g.Transactions, "GENESIS_BLOCK")
}	