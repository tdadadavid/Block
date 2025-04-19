package chain

import (
	"context"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/tdadadavid/block/pkg/block"
)

// LastKey tracks the block at "LAST" position
var LastKey = []byte("LAST")

type Storage interface {
	FindLastOrCreate(ctx context.Context) block.Block
	Create(ctx context.Context, key string, b block.Block) error
	FindByHash(ctx context.Context, hash string) (block.Block, error)
	FindLast(ctx context.Context) (block.Block, error)
	UpdateLast(ctx context.Context, b block.Block) error
}

type ChainStore struct {
	store  *badger.DB
	logger *slog.Logger
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
//   - error: Returns the error during a block creation process
func (cs *ChainStore) Create(_ context.Context, key string, b block.Block) error {
	err := cs.store.Update(func(txn *badger.Txn) (err error) {
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

// FindByHash finds a block by the given hash
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
//   - err(error): Returns the error during search
func (cs *ChainStore) FindByHash(_ context.Context, hash string) (b block.Block, err error) {
	b = block.Block{}
	err = cs.store.View(func(txn *badger.Txn) error {
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

// FindLast finds a block by the given hash
//
// Process:
//   - Retrieves the block in bytes if it exists, else returns error
//   - Deserializes the bytes into a block
//
// Returns
//   - block: Returns the block found
//   - error: Returns the error during search
func (cs *ChainStore) FindLast(_ context.Context) (block.Block, error) {
	b := &block.Block{}
	err := cs.store.View(func(txn *badger.Txn) error {
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
//   - error: Returns the error during an update process
func (cs *ChainStore) UpdateLast(_ context.Context, b block.Block) (err error) {
	err = cs.store.Update(func(txn *badger.Txn) error {
		data, err := b.Serialize()
		if err != nil {
			return err
		}
		err = txn.Set(LastKey, data)
		return err
	})
	return err
}
