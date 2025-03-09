package blockchain

import "github.com/tdadadavid/block/pkg/block"

type Blockchain struct {
	blocks []block.Block
}

func New() (bc Blockchain) {
	bc = Blockchain{
		blocks: []block.Block{block.NewGenesisBlock()},
	}
	return bc
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]                                     // get previous block
	newBlock := block.NewBlock(data, prevBlock.GetHash(), block.HASH_DIFFICULTY) // create new block
	bc.blocks = append(bc.blocks, newBlock)                                      // add block to the chain
}
