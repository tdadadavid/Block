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

func printChain() {
	blockChain.PrettyPrint()
}

func addBlock(data string) {
	blockChain.AddBlock(data)
}

func printLastBlockOnChain() {
	b, _ := blockChain.FindLast()
	fmt.Printf("Block {%v}\n", b)
}
