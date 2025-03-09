package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

var (
	HASH_DIFFICULTY uint = 4
)

type Block struct {
	// Timestamp captures the time the 'Block' was created
	Timestamp int64

	// Transactions
	Transactions string

	// PrevBlockHash is the Hash of the previous block
	PrevBlockHash string

	// Hash captures the Hash of the current block
	Hash string

	// Height is the history of blockchain
	Height uint

	// Nonce this is used in blockchain mining
	Nonce int32
}

// NewBlock creates & returns a new block on the blockchain
// It also runs the proofOfWork algorithm on the newly created block
func NewBlock(data string, prevBlkHash string, height uint) (block Block) {
	timestamp := time.Now()
	block = Block{
		Timestamp:     timestamp.Unix(),
		Transactions:  data,
		PrevBlockHash: prevBlkHash,
		Hash:          hex.EncodeToString(sha256.New().Sum([]byte(data))),
		Height:        height,
		Nonce:         0,
	}
	block.runProofOfWork()
	return block
}

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

func (b *Block) GetTransaction() string {
	return b.Transactions
}

func (b *Block) GetHash() string {
	return b.Hash
}

func (b *Block) runProofOfWork() {
	// first validate the Hash if it is not correct increment the Nonce to affect the Hash shuffling
	for !b.validate() {
		b.Nonce += 1
	}
	b.calculateHash()
}

// calculateHash calculates and updates the block's Hash value
func (b *Block) calculateHash() {
	data, err := b.prepareHashData()
	if err != nil {
		fmt.Printf("error preparing Hash data %v", data)
	}

	hash := sha256.Sum256(data)
	b.Hash = hex.EncodeToString(hash[:])
}

// validate validates the block
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
	hashPrefix := hexHash[:HASH_DIFFICULTY]

	// Compare with HASH_DIFFICULTY zeroes
	zeroString := strings.Repeat("0", int(HASH_DIFFICULTY))

	// compare the values to see if its valid
	valid = zeroString == hashPrefix // NOTE: after testing this, this computation ran approx 200K times for sample data

	return valid
}

func (b *Block) prepareHashData() ([]byte, error) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(b.clone())
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), err
}

func (b *Block) clone() *Block {
	return &Block{
		Timestamp:     b.Timestamp,
		Transactions:  b.Transactions,
		PrevBlockHash: b.PrevBlockHash,
		Hash:          b.Hash,
		Height:        HASH_DIFFICULTY,
		Nonce:         b.Nonce,
	}
}
