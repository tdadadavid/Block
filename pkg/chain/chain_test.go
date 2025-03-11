package chain

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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

	bc := New("./../../data/blocks")
	assert.NotNil(t, bc)
}

func TestBlockchain_AddBlock(t *testing.T) {
	defer cleanUp(t)

	bc := New("./../../data/blocks")

	assert.NotNil(t, bc)
	bc.AddBlock("test_data")

	iter := bc.iter()
	for iter.HasNext() {
		it := iter.Next()
		if it == nil {
			assert.Equal(t, "", iter.currentHash)
		} else {
			assert.NotEmpty(t, it.GetTransaction())
			assert.Equal(t, iter.currentHash, it.GetTransaction())
		}
	}
}
