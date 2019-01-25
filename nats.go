// nats.go
package main

import (
	"encoding/json"
	"log"

	nats "github.com/nats-io/go-nats"
)

var (
	nc   *nats.Conn
	subs *nats.Subscription
)

func natsConnect() {
	var err error

	nc, err = nats.Connect(servers)
	if err != nil {
		log.Panicln(err)
	}
}

func sendID() {
	msg, _ := makeMsg(UPDATE, []byte(myID), []byte{})
	nc.Publish(ServerChannel, msg)
}

func handleClientMsg(msg *nats.Msg) {
	mymsg := myMsg{}
	json.Unmarshal(msg.Data, &mymsg)
	switch mymsg.Cmd {
	case UPDATE:
		newchain := Chain{}
		json.Unmarshal(mymsg.Data, &newchain)
		chain.replaceChain(newchain)
		printBlockchain(ui)
		log.Printf("Update: %x\n", newchain.lastBlock().Hash[:4])
	default:
	}
}
