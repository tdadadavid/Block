package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/tdadadavid/block/pkg/block"
)


func (c *Chain) FindLast() (block.Block, error) {
	return c.store.FindLast(c.chainCtx)
}

func (c *Chain) PrintBlock(hash string) {
	block, err := c.store.FindByHash(context.Background(), hash)
	if err != nil {
		fmt.Printf("err printing block")
	}
	fmt.Println(block)
}

// PrintChain prints the blockchain in a formatted way
func (c *Chain) PrintChain() {
	var builder strings.Builder
	builder.WriteString("Blockchain {\n")
	builder.WriteString(fmt.Sprintf("  currentHash: %q\n", c.currentHash))
	
	// Get blocks from the store
	blocks, err := c.GetAllBlocks()
	if err != nil {
		fmt.Printf("Blockchain {\n  currentHash: %q\n  blocks: [Error fetching blocks: %v]\n}", 
			c.currentHash, err)
	}
	
	builder.WriteString("  blocks: [\n")
	for i, block := range blocks {
		builder.WriteString(fmt.Sprintf("    %s", block.String()))
		if i < len(blocks)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString("  ]\n")
	builder.WriteString("}")
	
	fmt.Println(builder.String())
}
