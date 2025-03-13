package chain

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanUp(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("./../../data")
		if err != nil {
			panic(err)
		}
	})
}

func TestBlockchain_New(t *testing.T) {
	defer cleanUp(t)

	bc := New(context.Background(), "./../../data/blocks")
	assert.NotNil(t, bc)
}

func TestBlockchain_AddBlock(t *testing.T) {
	defer cleanUp(t)

	// create blockchain
	bc := New(context.Background(), "./../../data/blocks")
	assert.NotNil(t, bc)

	// get prevHash (genesis block)
	prevHash := bc.getLastHash()
	assert.NotNil(t, prevHash)

	// add block
	bc.AddBlock("test_data")

	ctx := context.Background()
	
	iter := bc.iter()
	for iter.HasNext(ctx) {
		it := iter.Next(ctx)
		if it == nil {
			assert.Equal(t, "", iter.currentHash)
		} else {
			assert.NotEmpty(t, it.GetTransaction())
			// it the iteration moves backwards towards the genesis block
			assert.Equal(t, iter.currentHash, prevHash)
		}
	}
}
