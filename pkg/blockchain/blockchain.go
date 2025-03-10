package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/tdadadavid/block/pkg/block"
)

type Blockchain struct {
	// store In memory database that stores the blockchain's data
	store *badger.DB

	// currentHash the current hash for this blockchain
	currentHash string
}

// New creates a new blockchain
//
// Parameters:
//   - storagePath(string): The path to the in-memory storage
//
// Process:
//   - Opens or creates the storage location
//   - Creates the blockchain with the storage
//   - Assigns the currentHash of the chain to the LastHash on the block (in this case it will be the GenesisBlock)
//
// Notes:
//   - If the store fails to open or create this function panics
//
// Returns:
//   - bc(Blockchain): The newly created blockchain
func New(storagePath string) (bc Blockchain) {
	store, err := badger.Open(badger.DefaultOptions("./../../data/blocks"))
	if err != nil {
		panic(fmt.Errorf("failed to create blockchain %v", err))
	}

	bc = Blockchain{store: store}
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
func (bc *Blockchain) getLastHash() string {
	b := bc.findLastOrCreate()
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
func (bc *Blockchain) AddBlock(data string) {
	// get previous block
	prevBlock, err := bc.findLast()
	if err != nil {
		fmt.Printf("error while finding previous block: %s\n", err.Error())
	}

	// creates new block with previous block hash
	newBlock := block.NewBlock(data, prevBlock.GetHash(), block.HashDifficulty) // create new block
	err = bc.Create(newBlock.GetHash(), newBlock)
	if err != nil {
		fmt.Printf("error while creating new block %v", err)
	}

	// update 'LAST' key in blockchain point to new block.
	err = bc.UpdateLast(newBlock)
	if err != nil {
		fmt.Println("error while updating last block")
	}

	// update the blockchain current-hash
	bc.currentHash = newBlock.GetHash()
}

// iter add a block to the chain
//
// Process:
//   - Creates an iterator object with the currentHash of the chain and the blockchain itself
//
// Returns:
//   - iterator: The BlockchainIterator object
func (bc *Blockchain) iter() BlockChainIterator {
	return BlockChainIterator{
		currentHash: bc.currentHash,
		blockchain:  bc,
	}
}
