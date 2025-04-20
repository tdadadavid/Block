package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/tdadadavid/block/pkg/block"
)

// LastKey tracks the block at the "LAST" position
var LastKey = []byte("LAST")

type Storage interface {
	CreateBlock(ctx context.Context, key string, b block.Block) error
	FindBlockByHash(ctx context.Context, hash string) (block.Block, error)
	FindLastBlock(ctx context.Context) (block.Block, error)
	UpdateLastBlock(ctx context.Context, b block.Block) error
	FindAllWallets(ctx context.Context) ([][]byte, error)
}

type Store struct {
	store  *badger.DB
	logger *slog.Logger
}

func Open(path string) (s Storage, err error) {
	// Set the in-memory store for the chain and disable storage logs
	options := badger.DefaultOptions(path).WithLogger(nil)

	store, err := badger.Open(options)
	if err != nil {
		panic(fmt.Errorf("failed to create chain %v", err))
	}

	ss := Store{store: store, logger: slog.Default()}

	return &ss, err
}

func (s *Store) FindAllWallets(ctx context.Context) (wa [][]byte, err error) {
	err = s.store.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.Reverse = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			return item.Value(func(val []byte) (err error) {
				wa = append(wa, val)
				return err
			})
		}
		return err
	})
	return wa, err
}

// CreateBlock this creates new block on the chain
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
func (s *Store) CreateBlock(_ context.Context, key string, b block.Block) error {
	err := s.store.Update(func(txn *badger.Txn) (err error) {
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

// FindBlockByHash finds a block by the given hash
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
func (s *Store) FindBlockByHash(_ context.Context, hash string) (b block.Block, err error) {
	b = block.Block{}
	err = s.store.View(func(txn *badger.Txn) error {
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

// FindLastBlock finds a block by the given hash
//
// Process:
//   - Retrieves the block in bytes if it exists, else returns error
//   - Deserializes the bytes into a block
//
// Returns
//   - block: Returns the block found
//   - error: Returns the error during search
func (s *Store) FindLastBlock(_ context.Context) (block.Block, error) {
	b := &block.Block{}
	err := s.store.View(func(txn *badger.Txn) error {
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

// UpdateLastBlock This updates a special key on our chain called "LAST", this key stores the last block on the chain
//
// Process:
//   - Serializes the block
//   - Inserts the LastKey ("LAST") & serialized block in bytes into the storage
//
// Returns
//   - error: Returns the error during an update process
func (s *Store) UpdateLastBlock(_ context.Context, b block.Block) (err error) {
	err = s.store.Update(func(txn *badger.Txn) error {
		data, err := b.Serialize()
		if err != nil {
			return err
		}
		err = txn.Set(LastKey, data)
		return err
	})
	return err
}
