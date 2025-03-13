package cmd

import (
	"context"
	"fmt"

	"github.com/tdadadavid/block/pkg/chain"
)

var blockChain chain.Chain

func init() {
	blockChain = chain.New(context.Background(), "/data/blocks")
}

func addBlock(data string) {
	blockChain.AddBlock(data)
}

func printBlock(hash string) {
	blockChain.PrintBlock(hash)
} 

func printChain(_ string) {
	blockChain.PrintChain()
}

func printLastBlockOnChain(_ string) {
	b, _ := blockChain.FindLast()
	fmt.Printf("Block {%v}\n", b)
}
