// main.go
package main

import (
	"fmt"
)

func main() {
	genesis := Genesis()

	fmt.Printf("Genesis: %+v\n", genesis)

	data := Data{
		"test": "test",
	}

	block := MineBlock(genesis, data)

	fmt.Printf("%+v\n", block)

	chain := NewChain()
	fmt.Printf("chain: %+v\n", chain)

	chain.addBlock(data)
	fmt.Printf("newChain: %+v\n", chain)
}
