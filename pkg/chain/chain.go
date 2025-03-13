package chain

import (
	"context"
	"encoding/json"
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

// PrettyPrint this formats the chain into a readable json format for easy debugging
func (c *Chain) PrettyPrint() {
	json, err := json.MarshalIndent(c, "Blockchain ", "")
	if err != nil {
		fmt.Printf("error marshalling block: %s", err)
	}
	fmt.Println(string(json))
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

// Utility functions
func (c *Chain) FindLast() (block.Block, error) {
	return c.store.FindLast(c.chainCtx)
}
