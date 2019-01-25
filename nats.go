// nats.go
package main

import (
	"encoding/json"
	"log"
	"time"

	nats "github.com/nats-io/go-nats"
)

var (
	nc     *nats.Conn
	subs   *nats.Subscription
	server *nats.Subscription
)

func natsConnect() {
	var err error

	nc, err = nats.Connect(servers)
	if err != nil {
		log.Panicln(err)
	}
}

func sendID() {
	if MinerMode {

	} else if !ServerMode {

	}
}

func handleServerMsg(msg *nats.Msg) {
	mymsg := myMsg{}
	json.Unmarshal(msg.Data, &mymsg)
	switch mymsg.Cmd {
	case ADD:
		// SEND Data to miner
		last, err := json.Marshal(chain.lastBlock())
		if err != nil {
			log.Printf("Add err: %+v\n", err)
			return
		}
		data := mymsg.Data
		newmsg, err := makeMsg(MINE, data, last)
		if err != nil {
			log.Printf("Add err: %+v\n", err)
			return
		}
		// Send and wait the first response
		ret, err := nc.Request(MinerChannel, newmsg, time.Duration(MineRate*2)*time.Millisecond)

		// parse reply from miner
		newBlockMsg := myMsg{}
		json.Unmarshal(ret.Data, &newBlockMsg)

		newBlock := Block{}
		json.Unmarshal(newBlockMsg.Data, &newBlock)

		*chain = append(*chain, newBlock)
		log.Printf("New Chain: %xb - %s\n", chain.lastBlock().Hash[:8], string(newBlockMsg.Msg))
		isValid := isValidChain(chain)

		if !isValid {
			log.Printf("invalid chain: %+v\n", newBlock)
			return
		}

		// send to the miner the BLOCK_FOUND msg
		foundMsg, _ := makeMsg(BLOCK_FOUND, []byte{}, []byte{})
		nc.Publish(MinerChannel, foundMsg)

		jdata, err := json.Marshal(chain)
		if err != nil {
			log.Printf("BLOCKS ERR: %+v\n", err)
		}

		// update all clients
		upChain, _ := makeMsg(UPDATE, jdata, []byte{})
		nc.Publish(ClientChannel, upChain)
		nc.Publish(MinerChannel, upChain)

		printBlockchain(ui)

		log.Printf("Found Block: %s\n", string(newBlockMsg.Msg))

	default:
		log.Printf("Invalid CMD: %+v\n", mymsg)
		return
	}
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
	default:
	}
}

func handleMinerMsg(msg *nats.Msg) {
	mymsg := myMsg{}
	json.Unmarshal(msg.Data, &mymsg)
	switch mymsg.Cmd {
	case MINE:
		setNet(true)
		defer setNet(false)
		Found = false
		last := Block{}
		json.Unmarshal(mymsg.Msg, &last)
		lasthash := last.LastHash
		minedBlock, err := MineBlock(last, mymsg.Data)
		if err != nil {
			log.Printf("miner: %+v\n", err)
			return
		}
		retval, err := json.Marshal(minedBlock)
		if err != nil {
			log.Printf("new block err: %+v\n", err)
			return
		}

		newmsg, err := makeMsg(BLOCK_FOUND, retval, []byte(myID))

		log.Printf("sending block: %+v\n", lasthash[:8])
		nc.Publish(msg.Reply, newmsg)
	case BLOCK_FOUND:
		Found = true
	case UPDATE:
		newchain := Chain{}
		json.Unmarshal(mymsg.Data, &newchain)
		chain.replaceChain(newchain)
		printBlockchain(ui)
	default:
	}
}
