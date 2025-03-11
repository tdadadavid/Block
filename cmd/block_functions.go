package cmd

import (
	"fmt"

	"github.com/tdadadavid/block/pkg/chain"
)

var blockChain chain.Chain

func init() {
	blockChain = chain.New("/data/blocks")
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
