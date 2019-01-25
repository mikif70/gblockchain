// server.go
package main

import (
	"encoding/json"
	"log"
	"time"

	nats "github.com/nats-io/go-nats"
)

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
		newmsg, err := makeMsg(MINE, mymsg.Data, last)
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
		log.Printf("New Chain: %x - %s\n", newBlock.Hash[:4], string(newBlockMsg.Msg))

		if !isValidChain(chain) {
			log.Printf("invalid chain: %+v\n", newBlock)
			return
		}

		// send to the miner the BLOCK_FOUND msg
		foundMsg, _ := makeMsg(BLOCK_FOUND, newBlock.LastHash, []byte{})
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

		log.Printf("Found Block: %s - %x\n", string(newBlockMsg.Msg), newBlock.LastHash[:8])
	case UPDATE:
		jdata, err := json.Marshal(chain)
		if err != nil {
			log.Printf("BLOCKS ERR: %+v\n", err)
		}

		// update all clients
		upChain, _ := makeMsg(UPDATE, jdata, []byte{})
		nc.Publish(ClientChannel, upChain)
		nc.Publish(MinerChannel, upChain)
	default:
		log.Printf("Invalid CMD: %+v\n", mymsg)
		return
	}
}
