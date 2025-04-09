package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/tdadadavid/block/pkg/toolkit"
	"github.com/tdadadavid/block/pkg/transactions"
)

var (
	HashDifficulty int32 = 4
)

// A grouping of transactions, marked with a timestamp, and a
// fingerprint of the previous block. The block header is hashed to produce a proof of work,
// thereby validating the transactions. Valid blocks are added to the main blockchain by network consensus.
// Ref: https://cypherpunks-core.github.io/bitcoinbook/glossary.html
type Block struct {
	// Timestamp captures the time the 'Block' was created
	Timestamp int64 `json:"timestamp"`

	// Transactions
	Transactions []transactions.Transaction `json:"txns"`

	// PrevBlockHash is the Hash of the previous block
	PrevBlockHash string `json:"pbh"`

	// Hash captures the Hash of the current block
	Hash string `json:"hash"`

	// Height is the history of chain
	Height int32 `json:"height"`

	// Nonce this is used in chain mining
	Nonce int32 `json:"nonce"`

	logger *slog.Logger
}

// New creates & returns a new block on the chain
//
// Parameters:
//   - data(string): The transactional data on the chain
//   - prevBlkHash(string): The hash of the previous block
//   - height(int32): The height of the block
//
// Process:
//   - Creates the block with the given parameters and the current time in milliseconds
//   - It deserializes the transaction and hashes it, then converts to hexadecimal
//   - It also runs the proofOfWork algorithm on the newly created block
//
// Returns:
//   - block: The block that just got created
func New(data transactions.Transaction, prevBlkHash string, height int32) (block Block) {
	timestamp := time.Now()

	bytez, err := data.Serialize()
	if err != nil {
		panic(err)
	}
	val := sha256.Sum256(bytez)
	hash := hex.EncodeToString(val[:])

	block = Block{
		Timestamp:     timestamp.Unix(),
		Transactions:  []transactions.Transaction{data},
		PrevBlockHash: prevBlkHash,
		Hash:          hash,
		Height:        height,
		Nonce:         0,
		logger:        slog.Default(),
	}
	block.mine() // mine the block
	return block
}

// NewGenesisBlock creates the first block on the chain known as the 'GENESIS_BLOCK'
//
// Note:
//   - This method should be called once and that is during the BlockChain creation.
//   - Coinbase is the first coin (base) for cryptocurrency like bitcoin, it has no inputs
//
// Returns:
//   - genesis: the genesis block on the chains
func NewGenesisBlock(coinbase transactions.Transaction) (genesis Block) {
	timestamp := time.Now()
	genesis = Block{
		Timestamp:     timestamp.Unix(),
		Transactions:  []transactions.Transaction{coinbase},
		PrevBlockHash: "",
		Hash:          hex.EncodeToString(sha256.New().Sum([]byte(""))),
		Height:        0,
		Nonce:         0,
	}
	genesis.mine()
	return genesis
}

// GetTransaction returns the transaction of the block
func (b *Block) GetTransaction() []transactions.Transaction {
	return b.Transactions
}

// GetHash returns the hash of the block
func (b *Block) GetHash() string {
	return b.Hash
}

// GetPrevBlockHash returns the hash of the previous block
func (b *Block) GetPrevBlockHash() string {
	return b.PrevBlockHash
}

// GetHeight returns the height of the block
func (b *Block) GetHeight() int32 {
	return b.Height
}

// GetNonce returns the nonce of the block
func (b *Block) GetNonce() int32 {
	return b.Nonce
}

// GetTimestamp returns the timestamp of the block
func (b *Block) GetTimestamp() int64 {
	return b.Timestamp
}

// String returns a string representation of a Block
// This implements the Stringer interface to enable us printing like this fmt.Println(&block)
func (b *Block) String() string {
	return fmt.Sprintf("{Hash: %q, PrevBlockHash: %q, Height: %d, Timestamp: %d, Transactions: %q, Nonce: %d}",
		b.Hash, b.PrevBlockHash, b.Height, b.Timestamp, b.Transactions, b.Nonce)
}

// Serialize converts the Block in binary data stored in the storage
//
// Process:
//   - Writes the data types with known length into binary bytes (Height(int32), Nonce(int32), Timestamp(int64))
//   - Writes data types with varying length into binary bytes (Transactions(string), PrevHashBlock(string), and Hash(string))
//
// Returns:
//   - val: byte representation of block
//   - err: error during serialization process
func (b *Block) Serialize() (val []byte, err error) {
	var buf bytes.Buffer

	// Write Height
	if err = binary.Write(&buf, binary.LittleEndian, b.Height); err != nil {
		err = fmt.Errorf("error writing height: %w", err)
		return val, err
	}

	// Write Nonce
	if err = binary.Write(&buf, binary.LittleEndian, b.Nonce); err != nil {
		err = fmt.Errorf("error writing nonce: %w", err)
		return val, err
	}

	// Write Timestamp
	if err = binary.Write(&buf, binary.LittleEndian, b.Timestamp); err != nil {
		err = fmt.Errorf("error writing timestamp: %w", err)
		return val, err
	}

	// Write number of transactions
	txCount := uint32(len(b.Transactions))
	if err := binary.Write(&buf, binary.LittleEndian, txCount); err != nil {
		err = fmt.Errorf("error writing transaction count: %w", err)
		return val, err
	}

	// Write Transactions (each one serialized separately)
	for _, tx := range b.Transactions {
		txBytes, err := tx.Serialize() // Assuming Transaction has a Serialize method
		if err != nil {
			err = fmt.Errorf("error serializing transaction: %w", err)
			return val, err
		}

		txLen := uint32(len(txBytes))
		if err := binary.Write(&buf, binary.LittleEndian, txLen); err != nil {
			err = fmt.Errorf("error writing transaction length: %w", err)
			return val, err
		}
		buf.Write(txBytes)
	}

	// Write Hash
	if err := toolkit.SerializeString(&buf, b.Hash); err != nil {
		return val, err
	}

	// Write Previous Block Hash
	if err := toolkit.SerializeString(&buf, b.PrevBlockHash); err != nil {
		return val, err
	}

	return buf.Bytes(), err
}

// Deserialize converts the bytes data from the storage into a Block
//
// Parameters:
//   - data([]byte): the data from the storage from which the block object will be populated
//
// Process:
//   - Check if the block is empty and returns error if that is true
//   - Reads all data types with known length like (Height(int32), Nonce(int32), Timestamp(int64))
//   - Reads all varying data types eg (Transactions(string), PrevHashBlock(string), and Hash(string))
//
// Returns:
//   - err(error): error during deserialization process
func (b *Block) Deserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)

	// Read Height
	if err := binary.Read(buf, binary.LittleEndian, &b.Height); err != nil {
		err = fmt.Errorf("error reading height: %w", err)
		return err
	}

	// Read Nonce
	if err := binary.Read(buf, binary.LittleEndian, &b.Nonce); err != nil {
		err = fmt.Errorf("error reading nonce: %w", err)
		return err
	}

	// Read Timestamp
	if err := binary.Read(buf, binary.LittleEndian, &b.Timestamp); err != nil {
		err = fmt.Errorf("error reading timestamp: %w", err)
		return err
	}

	// Read number of transactions
	var txCount uint32
	if err := binary.Read(buf, binary.LittleEndian, &txCount); err != nil {
		err = fmt.Errorf("error reading transaction count: %w", err)
		return err
	}

	// Read Transactions
	b.Transactions = make([]transactions.Transaction, txCount)
	for i := uint32(0); i < txCount; i++ {
		// Read transaction length
		var txLen uint32
		if err = binary.Read(buf, binary.LittleEndian, &txLen); err != nil {
			err = fmt.Errorf("error reading transaction length: %w", err)
			return err
		}

		// Validate transaction length
		// ✅ Ensure txLen is within buffer range
		bufLen := uint32(buf.Len())
		if txLen == 0 || txLen > bufLen {
			err = fmt.Errorf("invalid transaction length: %d", txLen)
			return err
		}

		// Read transaction bytes safely
		txBytes := make([]byte, txLen)
		// ✅ Use io.ReadFull for safety
		if _, err = io.ReadFull(buf, txBytes); err != nil {
			err = fmt.Errorf("error reading transaction data: %w", err)
			return err
		}

		// Deserialize transaction
		var tx transactions.Transaction
		if err = tx.Deserialize(txBytes); err != nil {
			err = fmt.Errorf("error deserializing transaction: %w", err)
			return err
		}

		b.Transactions[i] = tx
	}

	// Read Hash
	hash, err := toolkit.DeserializeString(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		err = fmt.Errorf("error reading hash: %w", err)
		return err
	}
	b.Hash = hash

	// Read Previous Block Hash
	prevHash, err := toolkit.DeserializeString(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		err = fmt.Errorf("error reading previous block hash: %w", err)
		return err
	}
	b.PrevBlockHash = prevHash

	// if we get to the end of the input then
	if err == io.EOF {
		err = nil
	}
	return err
}

// mine validates and calculate the hash of the new block to be added
//
// Process:
//   - Validate the block, if the block is not validated increase the block Nonce to increase Hash shuffling
//   - calculateHash updates the hash of the current block
func (b *Block) mine() {
	// first validate the Hash if it is not correct increment the Nonce to affect the Hash shuffling
	for !b.validate() {
		b.Nonce += 1
	}
}

// validate validates the block
//
// Process:
//   - Get the binary representation of the current block
//   - Using SHA256 hash the bytes data, converts it to hexadecimal, then get the HashDifficulty prefix
//   - Generate a zeroString HashDifficulty with length
//   - compare hash and the zeroString for validity
//
// Note:
//   - After testing this, this computation ran approx 200K times for sample data
//   - This is the proof-of-work algorithm
//
// Returns:
//   - valid: true if the validation block passes else false.
func (b *Block) validate() (valid bool) {
	data, err := b.Serialize()
	if err != nil {
		return valid
	}

	// Hash the bytes of the cloned block and convert it to hexadecimal
	// compute hash for block
	hash := sha256.Sum256(data)
	// get hexadecimal value of the hash
	hexHash := hex.EncodeToString(hash[:])
	// get the HASH_DIFFICULTY prefix from hash
	hashPrefix := hexHash[:HashDifficulty]

	// Compare with HASH_DIFFICULTY zeroes
	zeroString := strings.Repeat("0", int(HashDifficulty))

	// compare the values to see if its valid
	valid = zeroString == hashPrefix

	// set the Hash of the block if it is valid
	if valid {
		b.Hash = hexHash
	}

	return valid
}
