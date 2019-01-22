// main.go
package main

import (
	"fmt"
	"log"
	"os"
	//	"time"
	//	ui "github.com/jroimartin/gocui"
	//	ui "github.com/gizak/termui"
)

var (
	chain *Chain
)

func main() {

	fs, err := os.OpenFile("./chain.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Log file error: ", err.Error())
		os.Exit(-4)
	}
	defer fs.Close()

	log.SetOutput(fs)

	genesis := Genesis()

	fmt.Printf("Genesis: %+v\n", genesis)

	data := Data{
		"test": "test",
	}

	block := MineBlock(genesis, data)

	fmt.Printf("%+v\n", block)

	chain = NewChain()
	fmt.Printf("chain: %+v\n", chain)

	chain.addBlock(data)
	fmt.Printf("newChain: %+v\n", chain)

	RunConsole()
}
