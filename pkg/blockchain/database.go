package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/tdadadavid/block/pkg/block"
)

// LastKey tracks the block at "LAST" position
var LastKey = []byte("LAST")

// findLastOrCreate used to get the last block in the chain or creates a new genesis block
//
// Process:
//   - Queries storage to get the block at the "LAST" position
//   - If the block is not found then create the Genesis block for this chain, else returns the last found block
//
// Returns
//   - last: The block at the "LAST" position
func (bc *Blockchain) findLastOrCreate() (last block.Block) {
	// try to findLast the block with the key 'LAST'
	last, err := bc.findLast()
	if err != nil {
		//TODO: standardize logging
		fmt.Printf("error finding last block err: [%s]\n", err.Error())
	}

	fmt.Printf("found last item %s\n", last.Transactions)

	if last.Transactions == "" {
		// if not found create the genesis block
		fmt.Printf("creating new genesis block\n")
		last = block.NewGenesisBlock()
		err = bc.Create(string(LastKey), last)
		if err != nil {
			fmt.Printf("error while creating last item %s", err.Error())
		}
	}

	return last
}

// Create this creates new block on the chain
//
// Parameters:
//   - key(string): The hash of the block that serves as key for the block in the storage
//   - b(Block): The block to be stored in the storage
//
// Process:
//   - Serializes the block
//   - Inserts the key (block's hash) & serialized block in bytes into the storage
//
// Returns
//   - error: Returns the error during block creation process
func (bc *Blockchain) Create(key string, b block.Block) error {
	err := bc.store.Update(func(txn *badger.Txn) (err error) {
		data, err := b.Serialize()
		if err != nil {
			return err
		}
		err = txn.Set([]byte(key), data)
		if err != nil {
			return err
		}
		return err
	})
	return err
}

// findByHash finds a block by the given hash
//
// Parameters:
//   - hash(string): The hash of the block that we want to retrieve
//
// Process:
//   - Retrieves the block in bytes if it exists, else returns error
//   - Deserializes the bytes into a block
//
// Returns
//   - b(block): Returns the block just found
//   - err(error): Returns the error during block creation process
func (bc *Blockchain) findByHash(hash string) (b block.Block, err error) {
	b = block.Block{}
	err = bc.store.View(func(txn *badger.Txn) error {
		lastBlock, err := txn.Get([]byte(hash))
		if err != nil {
			return err
		}
		// deserialize the bytes into a Block struct
		return lastBlock.Value(func(val []byte) error {
			return b.Deserialize(val)
		})
	})

	return b, err
}

// findByHash finds a block by the given hash
//
// Process:
//   - Retrieves the block in bytes if it exists, else returns error
//   - Deserializes the bytes into a block
//
// Returns
//   - block: Returns the block just found
//   - error: Returns the error during block creation process
func (bc *Blockchain) findLast() (block.Block, error) {
	b := &block.Block{}
	err := bc.store.View(func(txn *badger.Txn) error {
		lastBlock, err := txn.Get(LastKey)
		if err != nil {
			return err
		}
		// deserialize the bytes into a Block struct
		return lastBlock.Value(func(val []byte) error {
			return b.Deserialize(val)
		})
	})

	return *b, err
}

// UpdateLast This updates a special key on our chain called "LAST", this key stores the last block on the chain
//
// Process:
//   - Serializes the block
//   - Inserts the LastKey ("LAST") & serialized block in bytes into the storage
//
// Returns
//   - error: Returns the error during update process
func (bc *Blockchain) UpdateLast(b block.Block) (err error) {
	err = bc.store.Update(func(txn *badger.Txn) error {
		data, err := b.Serialize()
		if err != nil {
			return err
		}
		err = txn.Set(LastKey, data)
		return err
	})
	return err
}
