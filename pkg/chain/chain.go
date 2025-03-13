package chain

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/tdadadavid/block/pkg/block"
)

type Chain struct {
	// store In memory database that stores the chain's data
	store *ChainStore

	// currentHash the current hash for this chain
	currentHash string

	chainCtx context.Context
}

// New creates a new chain
//
// Parameters:
//   - storagePath(string): The path to the in-memory storage
//
// Process:
//   - Opens or creates the storage location
//   - Creates the chain with the storage
//   - Assigns the currentHash of the chain to the LastHash on the block (in this case it will be the GenesisBlock)
//
// Notes:
//   - If the store fails to open or create this function panics
//
// Returns:
//   - bc(Chain): The newly created chain
func New(ctx context.Context, storagePath string) (bc Chain) {
	// Set the in-memory store for the chain and disable storage logs
	options := badger.DefaultOptions("./../../data/blocks").WithLogger(nil)

	store, err := badger.Open(options)
	if err != nil {
		panic(fmt.Errorf("failed to create chain %v", err))
	}

	bc = Chain{
		chainCtx: ctx,
		store:    &ChainStore{store: store},
	}
	bc.currentHash = bc.getLastHash()

	return bc
}

// getLastHash returns the hash of block as "LAST" position
//
// Process:
//   - finds block at last position or creates genesis block
//
// Returns:
//   - The hash of the block at "LAST" position
func (c *Chain) getLastHash() string {
	b := c.store.FindLastOrCreate(c.chainCtx)
	return b.GetHash()
}

// AddBlock add a block to the chain
//
// Parameters:
//   - data(string): The transactional data to be stored in the block
//
// Process:
//   - finds the previous block (block in the "LAST" position)
//   - Creates new block with given data and previous block's hash
//   - Updates the "LAST" key in the storage with the newly created block
//   - Set the Chains hash to the new block's hash
func (c *Chain) AddBlock(data string) {
	// get previous block
	prevBlock, err := c.store.FindLast(c.chainCtx)
	if err != nil {
		fmt.Printf("error while finding previous block: %s\n", err.Error())
	}

	// creates new block with previous block hash
	newBlock := block.New(data, prevBlock.GetHash(), block.HashDifficulty) // create new block
	err = c.store.Create(c.chainCtx, newBlock.GetHash(), newBlock)
	if err != nil {
		fmt.Printf("error while creating new block %v", err)
	}

	// update 'LAST' key in chain point to new block.
	err = c.store.UpdateLast(c.chainCtx, newBlock)
	if err != nil {
		fmt.Println("error while updating last block")
	}

	// update the chain current-hash
	c.currentHash = newBlock.GetHash()
}

// GetAllBlocks retrieves all blocks from the chain store
//
// Process:
// 	- It first creates an iterator object for the loop process
//  - while there is a block on the chain it gets it and appends it to a slice of blocks
//  - Reverse the blocks to get them in chronological order (oldest first)
//
// Returns:
// 	- `blocks []*block.Block`: a slice of blocks for this chain
//  - `err error`: an error object
func (c *Chain) GetAllBlocks() (blocks []*block.Block, err error) {
	iter := c.iter()

	for iter.HasNext(c.chainCtx) {
		block := iter.Next(c.chainCtx)
		if block == nil {
			err = fmt.Errorf("failed to get block %s", c.currentHash)
			return blocks, err
		}
		blocks = append(blocks, block)
	}

	// Reverse the blocks to get them in chronological order (oldest first)
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	
	return blocks, err
}


// iter add a block to the chain
//
// Process:
//   - Creates an iterator object with the currentHash of the chain and the chain itself
//
// Returns:
//   - iterator: The BlockchainIterator object
func (c *Chain) iter() ChainIterator {
	return ChainIterator{
		currentHash: c.currentHash,
		blockchain:  c,
	}
}


