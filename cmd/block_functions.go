package cmd

import (
	"context"
	"fmt"

	"github.com/tdadadavid/block/pkg/chain"
	"github.com/tdadadavid/block/pkg/transactions"
)

var blockChain chain.Chain
var chainStorePath = "/data/blocks"

func init() {
	blockChain = chain.New(context.Background(), chainStorePath)
}

func addBlock(data string) {
	blockChain.AddBlock(transactions.Transaction{}) //TODO: fix me
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
