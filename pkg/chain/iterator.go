package chain

import "github.com/tdadadavid/block/pkg/block"

// Iterator captures the building block of any iterable
type Iterator interface {
	Next() *block.Block
	HasNext() bool
}

// ChainIterator Iterator used for moving through the chain it implements the Iterator interface
type ChainIterator struct {
	// currentHash holds the hash of the currentBlock we are on in the chain
	currentHash string

	// blockchain holds a pointer to the Chain itself
	blockchain *Chain
}

// HasNext control the iteration over the chain
//
// Process:
//   - if the chain is nil no need to iterate stop by returning false
//   - if the currentHash is empty "" then we know we've reached the end of the chain
//   - we check if the currentHash points to a valid block, if it does we continue iteration else stop
//
// Returns:
//   - bool: either true or false to signify if the iteration should continue
func (it *ChainIterator) HasNext() bool {
	// if the chain is empty or the current hash is empty stop iteration
	if it.blockchain == nil || it.currentHash == "" {
		return false
	}

	// check if the block of the currentHash exists if it does the continue iteration else stop.
	_, err := it.blockchain.findByHash(it.currentHash)
	if err != nil {
		return false
	}

	return true
}

// Next moves to the next block on the chain if there is else returns a nil pointer
//
// Process:
//   - Check if there is a next using the HasNext method if it does continue iteration else returns nil and set the currentHash to empty.
//   - Find the current-block using the currentHash, if the block is not found return nil and set the currentHash to empty.
//   - Set the currentHash to the PreviousHash of the currentBlock, this way we move backwards to the genesis block
//
// Returns:
//   - currBlock: The pointer to the current block from the database
func (it *ChainIterator) Next() (curBlock *block.Block) {
	// check if there is a next block on the chain
	if it.HasNext() {
		it.currentHash = ""
		return curBlock
	}

	// find the current block using the current hash
	b, err := it.blockchain.findByHash(it.currentHash)
	if err != nil {
		it.currentHash = ""
		return curBlock
	}

	// set the currentHash to the PreviousBlock Hash (we move backwards)
	curBlock = &b
	it.currentHash = curBlock.PrevBlockHash

	return curBlock
}
