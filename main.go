// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	//	"time"
	// "github.com/jroimartin/gocui"
	//	ui "github.com/gizak/termui"
)

var (
	chain   *Chain
	myID    string
	clients map[string]bool
	miners  map[string]bool
)

func main() {

	flag.Parse()

	fs, err := os.OpenFile(LogFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Log file error: ", err.Error())
		os.Exit(-4)
	}
	defer fs.Close()

	log.SetOutput(fs)

	if ServerMode {
		data := "first block"
		chain = NewChain()
		chain.addBlock(data)
	}

	natsConnect()
	defer nc.Close()
	if ServerMode {
		subs, err = nc.QueueSubscribe(ServerChannel, ServerQueue, handleServerMsg)
	} else if MinerMode {
		myID, _ = newUUID()
		subs, err = nc.Subscribe(MinerChannel, handleMinerMsg)
		sendID()
	} else {
		subs, err = nc.Subscribe(ClientChannel, handleClientMsg)
		sendID()
	}
	defer subs.Unsubscribe()

	runConsole()
}
