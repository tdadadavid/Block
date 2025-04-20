package chain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tdadadavid/block/pkg/block"
	"github.com/tdadadavid/block/pkg/toolkit"
	"github.com/tdadadavid/block/pkg/transactions"

	"github.com/dgraph-io/badger/v4"
)

type Chain struct {
	// store In memory database that stores the chain's data
	store *ChainStore

	// currentHash the current hash for this chain
	currentHash string

	// chainCtx is the context for the chain, it is used to control the execution of the chain
	chainCtx context.Context

	logger *slog.Logger
}

// New instantiates a new chain from the store
//
// Parameters:
//   - `storagePath(string)`: The path to the in-memory storage
//
// Process:
//   - Opens or creates the storage location
//   - Creates the chain with the storage
//   - Assigns the currentHash of the chain to the LastHash on the block (in this case it will be the GenesisBlock)
//
// Notes:
//   - If the store fails to open or create this function panics
//   - This is called for chains that already exists
//
// Returns:
//   - `bc(Chain)`: The newly created chain
func New(ctx context.Context, storagePath string) (bc Chain) {
	// Set the in-memory store for the chain and disable storage logs
	options := badger.DefaultOptions(storagePath).WithLogger(nil)

	store, err := badger.Open(options)
	if err != nil {
		panic(fmt.Errorf("failed to create chain %v", err))
	}

	bc = Chain{
		chainCtx: ctx,
		store:    &ChainStore{store: store},
		logger:   slog.Default(),
	}

	// get the last block's hash
	hash, err := bc.getLastHash()
	if err != nil {
		panic(fmt.Errorf("failed to created chain %s", err))
	}

	bc.currentHash = hash

	return bc
}

// NewChain creates a new chain containing the coinbase transaction
//
// Parameters
//   - `ctx context.Context`: The context that control execution
//   - `name string`: The name of the blockchain. It is used when creating the store
//   - `address string`: The address of the blockchain
//
// Process
//   - The function tries to create a store if it fails then it panics
//     if it doesn't then it creates a coinbase transaction using the address & the genesis block
//     given and the coinbase data which is stored in the blockchain and in the 'LAST' position in the block
//
// Returns
//   - `bc Chain`: The newly created chain
func NewChain(ctx context.Context, name, address string) (bc Chain) {
	options := badger.DefaultOptions(fmt.Sprintf("./data/%s/blocks", name)).WithLogger(nil)

	store, err := badger.Open(options)
	if err != nil {
		panic(fmt.Errorf("failed to create chain %v", err))
	}

	// create coinbase transaction & genesis block
	cbtx := transactions.NewCoinbase(address, transactions.COINBASE_DATA)
	genesis := block.NewGenesisBlock(*cbtx)

	// create the chain
	bc = Chain{
		chainCtx: ctx,
		store:    &ChainStore{store: store},
		logger:   slog.Default(),
	}

	// store the genesis block in the store and in the 'LAST' position
	err = bc.store.Create(ctx, genesis.GetHash(), genesis)
	if err != nil {
		fmt.Printf("error while creating new block %v", err)
		return
	}

	// update 'LAST' key in chain point to genesis block.
	err = bc.store.UpdateLast(bc.chainCtx, genesis)
	if err != nil {
		fmt.Println("error while updating last block")
		return
	}

	// update the chain with the last hash
	bc.currentHash = genesis.GetHash()

	return bc
}

// FindUnspentTransactionsOutputs FindUnspentTransactions this get the total unspent transaction
//
// Parameters
//   - `ctx context.Context`: the context that controls the execution
//
// Process
//   - First it gets the unspent transactions for the user by going through every transaction (output) on each block
//     checking if the user can unlock that transaction (spend it) its then stored in a map,
//
// NOTE
//   - unspent transactions are the outputs (vouts) while the inputs are the spent transactions
//
// Returns
func (c *Chain) FindUnspentTransactionsOutputs(ctx context.Context) map[string]transactions.TxnOutputs {
	// tracks UTXO (unspent transaction outputs)
	utxos := make(map[string]transactions.TxnOutputs)

	// tracks the spent outputs for a transaction
	spentUTXOs := make(map[string][]int)

	iter := c.iter()

	for iter.HasNext(ctx) {
		curBlock := iter.Next(ctx)

		// get transaction for block
		for _, txn := range curBlock.GetTransaction() {

		OutputLoop:
			// get all outputs
			for outIdx, out := range txn.GetOutputs() {
				if spentUTXOs[txn.GetId()] != nil {
					// check if this transaction is in the spent transaction outputs
					for _, spentOutput := range spentUTXOs[txn.GetId()] {
						// if the current output has been spent, goto the outer loop & skip this inner
						if spentOutput == outIdx {
							continue OutputLoop
						}
					}
				}

				// add the output to the map of unspent outputs
				outs := utxos[txn.GetId()]
				outs.Outputs = append(outs.Outputs, out)
				utxos[txn.GetId()] = outs
			}

			if !txn.IsCoinbase() {
				// get all inputs for this transaction and mark them spent
				for _, in := range txn.GetInputs() {
					spentUTXOs[in.TxnId] = append(spentUTXOs[in.TxnId], int(in.Output))
				}
			}
		}
	}
	return utxos
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
func (c *Chain) AddBlock(data transactions.Transaction) {
	// get previous block
	prevBlock, err := c.store.FindLast(c.chainCtx)
	if err != nil || toolkit.Ref(prevBlock) == nil {
		fmt.Printf("error while finding previous block: %s\n", err)
		return
	}

	// creates new block with previous block hash
	newBlock := block.New(data, prevBlock.GetHash(), block.HashDifficulty) // create new block
	err = c.store.Create(c.chainCtx, newBlock.GetHash(), newBlock)
	if err != nil {
		fmt.Printf("error while creating new block %v", err)
		return
	}

	// update 'LAST' key in chain point to new block.
	err = c.store.UpdateLast(c.chainCtx, newBlock)
	if err != nil {
		fmt.Println("error while updating last block")
		return
	}

	// update the chain current-hash
	c.currentHash = newBlock.GetHash()
}

// GetAllBlocks retrieves all blocks from the chain store
//
// Process:
//   - It first creates an iterator object for the loop process
//   - while there is a block on the chain it gets it and appends it to a slice of blocks
//   - Reverse the blocks to get them in chronological order (oldest first)
//
// Returns:
//   - `blocks []*block.Block`: a slice of blocks for this chain
//   - `err error`: an error object
func (c *Chain) GetAllBlocks() (blocks []*block.Block, err error) {
	iter := c.iter()

	for iter.HasNext(c.chainCtx) {
		curBlock := iter.Next(c.chainCtx)
		if curBlock == nil {
			err = fmt.Errorf("failed to get curBlock %s", c.currentHash)
			return blocks, err
		}
		blocks = append(blocks, curBlock)
	}

	// Reverse the blocks to get them in chronological order (oldest first)
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	return blocks, err
}

// iter creates an iterator that iterates over the blockchain
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

// getLastHash returns the hash of block as "LAST" position
//
// Process:
//   - finds block at last position
//
// Returns:
//   - The hash of the block at "LAST" position
func (c *Chain) getLastHash() (val string, err error) {
	b, err := c.store.FindLast(c.chainCtx)
	if err != nil {
		c.logger.Error("no block in LAST position", slog.Any("error", err))
		err = fmt.Errorf("error retrieving LAST block %s", err)
		return val, err
	}
	val = b.GetHash()
	return val, err
}
