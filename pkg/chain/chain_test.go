package chain

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdadadavid/block/pkg/transactions"
)

func cleanUp(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("./data")
		if err != nil {
			panic(err)
		}
	})
}

func TestBlockchain_NewChain(t *testing.T) {
	defer cleanUp(t)

	bc := NewChain(context.Background(), "bitcoin", "0x0000")
	assert.NotNil(t, bc)
}

func TestBlockchain_AddBlock(t *testing.T) {
	defer cleanUp(t)

	// create blockchain with the coinbase (the initail coin release)
	bc := NewChain(context.Background(), "bitcoin", "0x000")
	assert.NotNil(t, bc)

	// add block
	bc.AddBlock(transactions.Transaction{})

	// get prevHash
	prevHash, _ := bc.getLastHash()
	assert.NotNil(t, prevHash)

	ctx := context.Background()

	iter := bc.iter()
	for iter.HasNext(ctx) {
		assert.NotEmpty(t, iter.currentHash) // as long as there is a block on the chain, there will always be a `currentHash`

		it := iter.Next(ctx)
		assert.NotEmpty(t, it.GetTransaction())
	}
}
