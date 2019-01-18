// main.go
package main

import (
	"fmt"
)

func main() {
	b := &Block{}

	genesis := Genesis()

	fmt.Printf("Genesis: %+v\n", genesis)

	data := &Data{
		"test": "test",
	}

	block := b.MineBlock(genesis, data)

	fmt.Printf("%+v\n", block)
}
