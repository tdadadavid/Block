package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

var (
	HashDifficulty int32 = 4
)

// Block a representation of blocks on chain
type Block struct {
	// Timestamp captures the time the 'Block' was created
	Timestamp int64 `json:"timestamp"`

	// Transactions
	Transactions string `json:"transactions"`

	// PrevBlockHash is the Hash of the previous block
	PrevBlockHash string `json:"previous_block_hash"`

	// Hash captures the Hash of the current block
	Hash string `json:"hash"`

	// Height is the history of chain
	Height int32 `json:"height"`

	// Nonce this is used in chain mining
	Nonce int32 `json:"nonce"`

	logger *slog.Logger
}

// NewBlock creates & returns a new block on the chain
//
// Parameters:
//   - data(string): The transactional data on the chain
//   - prevBlkHash(string): The hash of the previous block
//   - height(int32): The height of the block
//
// Process:
//   - Creates the block with the given parameters and the current time in milliseconds
//   - It also runs the proofOfWork algorithm on the newly created block
//
// Returns:
//   - block: The block that just got created
func NewBlock(data string, prevBlkHash string, height int32) (block Block) {
	timestamp := time.Now()
	block = Block{
		Timestamp:     timestamp.Unix(),
		Transactions:  data,
		PrevBlockHash: prevBlkHash,
		Hash:          hex.EncodeToString(sha256.New().Sum([]byte(data))),
		Height:        height,
		Nonce:         0,
		logger:        &slog.Logger{},
	}
	block.runProofOfWork()
	return block
}

// NewGenesisBlock creates the first block on the chain known as the 'GENESIS_BLOCK'
//
// Note:
//   - This method should be called once and that is during the BlockChain creation.
//
// Returns:
//   - genesis: the genesis block on the chains
func NewGenesisBlock() (genesis Block) {
	timestamp := time.Now()
	genesis = Block{
		Timestamp:     timestamp.Unix(),
		Transactions:  "GENESIS_BLOCK",
		PrevBlockHash: "",
		Hash:          hex.EncodeToString(sha256.New().Sum([]byte(""))),
		Height:        0,
		Nonce:         0,
	}
	return genesis
}

// GetTransaction returns the transaction of the block
func (b *Block) GetTransaction() string {
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
			return val, fmt.Errorf("error writing height: %w", err)
	}

	// Write Nonce
	if err = binary.Write(&buf, binary.LittleEndian, b.Nonce); err != nil {
		return val, fmt.Errorf("error writing nonce: %w", err)
	}

	// Write Timestamp
	if err = binary.Write(&buf, binary.LittleEndian, b.Timestamp); err != nil {
		return val, fmt.Errorf("error writing timestamp: %w", err)
	}

	// Write Transactions
	txBytes := []byte(b.Transactions)
	txLen := uint32(len(txBytes))
	if err = binary.Write(&buf, binary.LittleEndian, txLen); err != nil {
		return val, fmt.Errorf("error getting length of transaction string: %w", err)
	}
	buf.Write(txBytes)

	// Write Hash
	hashBytes := []byte(b.Hash)
	hashLen := uint32(len(hashBytes))
	if err = binary.Write(&buf, binary.LittleEndian, hashLen); err != nil {
		return val, fmt.Errorf("eerror getting length of hash string: %w", err)
	}
	buf.Write(hashBytes)

	// Write Previous Block Hash
	prevHashBytes := []byte(b.PrevBlockHash)
	prevHashLen := uint32(len(prevHashBytes))
	if err = binary.Write(&buf, binary.LittleEndian, prevHashLen); err != nil {
		return val, fmt.Errorf("eerror getting length of prevBlockHash string: %w", err)
	}
	buf.Write(prevHashBytes)

	val = buf.Bytes()
	return val, err
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
	if b == nil {
		return fmt.Errorf("block is nil")
	}

	buf := bytes.NewReader(data)

	// Read Height
	if err = binary.Read(buf, binary.LittleEndian, &b.Height); err != nil {
		return fmt.Errorf("error reading height: %w", err)
	}

	// Read Nonce
	if err = binary.Read(buf, binary.LittleEndian, &b.Nonce); err != nil {
		return fmt.Errorf("error reading nonce: %w", err)
	}

	// Read Timestamp
	if err = binary.Read(buf, binary.LittleEndian, &b.Timestamp); err != nil {
		return fmt.Errorf("error reading timestamp: %w", err)
	}

	// Read Transactions
	var txLen uint32
	if err = binary.Read(buf, binary.LittleEndian, &txLen); err != nil {
		return fmt.Errorf("error reading transaction length: %w", err)
	}

	if uint32(buf.Len()) < txLen {
		return fmt.Errorf("unexpected EOF while reading transactions (expected %d, found %d)", txLen, buf.Len())
	}

	txBytes := make([]byte, txLen)
	if _, err = buf.Read(txBytes); err != nil {
		return fmt.Errorf("error reading transaction data: %w", err)
	}
	b.Transactions = string(txBytes)

	// Read Hash
	var hashLen uint32
	if err = binary.Read(buf, binary.LittleEndian, &hashLen); err != nil {
		return fmt.Errorf("error reading hash length: %w", err)
	}

	if uint32(buf.Len()) < hashLen {
		return fmt.Errorf("unexpected EOF while reading hash (expected %d, found %d)", hashLen, buf.Len())
	}

	hashBytes := make([]byte, hashLen)
	if _, err = buf.Read(hashBytes); err != nil {
		return fmt.Errorf("error reading hash data: %w", err)
	}
	b.Hash = string(hashBytes)

	// Read Previous Block Hash
	var prevHashLen uint32
	if err = binary.Read(buf, binary.LittleEndian, &prevHashLen); err != nil {
		return fmt.Errorf("error reading previous hash length: %w", err)
	}

	if uint32(buf.Len()) < prevHashLen {
		return fmt.Errorf("unexpected EOF while reading previous hash (expected %d, found %d)", prevHashLen, buf.Len())
	}

	prevHashBytes := make([]byte, prevHashLen)
	if _, err = buf.Read(prevHashBytes); err != nil {
		return fmt.Errorf("error reading previous hash data: %w", err)
	}
	b.PrevBlockHash = string(prevHashBytes)

	return err
}

// runProofOfWork validates and calculate the hash of the new block to be added
//
// Process:
//   - Validate the block, if the block is not validated increase the block Nonce to increase Hash shuffling
//   - calculateHash updates the hash of the current block
func (b *Block) runProofOfWork() {
	// first validate the Hash if it is not correct increment the Nonce to affect the Hash shuffling
	for !b.validate() {
		b.Nonce += 1
	}
	b.calculateHash()
}

// calculateHash calculates and updates the block's Hash value
//
// Process:
//   - Gets the binary representation of the current block.
//   - Hash the bytes and encode it into hexadecimal
//   - Update the block Hash
func (b *Block) calculateHash() {
	data, err := b.prepareHashData()
	if err != nil {
		fmt.Printf("error preparing Hash data %v", data)
	}

	hash := sha256.Sum256(data)
	b.Hash = hex.EncodeToString(hash[:])
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
//
// Returns:
//   - valid: true if the validation block passes else false.
func (b *Block) validate() (valid bool) {
	data, err := b.prepareHashData()
	if err != nil {
		return false
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

	return valid
}

// prepareHashData serializes the blocks
//
// Returns:
//   - []bytes: the bytes representation of the block.
//   - error: the error that occurred while serialization.
func (b *Block) prepareHashData() ([]byte, error) {
	return b.Serialize()
}
