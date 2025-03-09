package blockchain

import "testing"

func TestBlockchain_New(t *testing.T) {

	bc := New()

	tests := []string{"data", "data1", "data2"}
	for _, tc := range tests {
		bc.AddBlock(tc)
	}

	// test for expected number of blocks
	if len(bc.blocks) != 4 {
		t.Errorf("failed creating blocks on blockchain,expected(4) got %d", len(bc.blocks))
	}

	// test that the first block is the Genesis
	genesis := bc.blocks[0]
	if genesis.GetTransaction() != "GENESIS_BLOCK" {
		t.Errorf("failed creating gensis block expected(GENESIS_BLOCK) got (%s)", genesis.GetTransaction())
	}

	// test the remaining blocks created.
	for idx := range 3 {
		transaction := bc.blocks[idx+1].GetTransaction()
		if transaction != tests[idx] {
			t.Errorf("failed creating transactions correctly expected(%s) got %s", tests[idx], transaction)
		}
	}
}
