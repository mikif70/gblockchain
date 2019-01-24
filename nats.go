// nats.go
package main

import (
	"encoding/json"
	"log"
	"time"

	nats "github.com/nats-io/go-nats"
)

var (
	nc   *nats.Conn
	subs *nats.Subscription
)

func NatsConnect() {
	var err error

	nc, err = nats.Connect(servers)
	if err != nil {
		log.Panicln(err)
	}

	//	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	//	if err != nil {
	//		nc.Close()
	//		log.Fatal(err)
	//	}
}

/*
func publish(cmd string, data string) {

	var jdata []byte
	var err error

	switch cmd {
	case "BLOCKS":
		jdata, err = json.Marshal(chain)
		if err != nil {
			log.Printf("BLOCKS ERR: %+v\n", err)
		}
	case "MINE":
	default:
	}

	if ServerMode {
		err = nc.Publish(ClientChannel, jdata)
	} else {
		err = nc.Publish(ServerChannel, jdata)
	}

	if err != nil {
		log.Printf("publish err: %+v\n", err)
		return
	}

	nc.Flush()
}
*/

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
		ret, err := nc.Request(MinerChannel, newmsg, time.Duration(MineRate*2)*time.Millisecond)

		// Parse reply from miner
		newBlockMsg := myMsg{}
		json.Unmarshal(ret.Data, &newBlockMsg)

		newBlock := Block{}
		json.Unmarshal(newBlockMsg.Data, &newBlock)

		*chain = append(*chain, newBlock)
		isValid := isValidChain(chain)

		if !isValid {
			log.Printf("invalid chain: %+v\n", newBlock)
			return
		}

		foundMsg, _ := makeMsg(BLOCK_FOUND, []byte{}, []byte{})
		nc.Publish(MinerChannel, foundMsg)

		jdata, err := json.Marshal(chain)
		if err != nil {
			log.Printf("BLOCKS ERR: %+v\n", err)
		}

		upChain, _ := makeMsg(UPDATE, jdata, []byte{})
		nc.Publish(ClientChannel, upChain)

		printBlockchain(ui)

		log.Printf("Found Block: %+v\n", newBlock)

	default:
		log.Printf("Invalid CMD: %+v\n", mymsg)
		return
	}
	//	log.Printf("Msg: %+v\n", mymsg)
}

func handleClientMsg(msg *nats.Msg) {
	log.Printf("Msg: %+v\n", msg)
}

func handleMinerMsg(msg *nats.Msg) {
	mymsg := myMsg{}
	json.Unmarshal(msg.Data, &mymsg)
	switch mymsg.Cmd {
	case MINE:
		Found = false
		last := Block{}
		json.Unmarshal(mymsg.Msg, &last)
		lasthash := last.LastHash
		minedBlock, err := MineBlock(last, mymsg.Data)
		if err != nil {
			log.Printf("miner: %+v\n", err)
		}
		retval, err := json.Marshal(minedBlock)
		if err != nil {
			log.Printf("new block err: %+v\n", err)
			return
		}

		newmsg, err := makeMsg(BLOCK_FOUND, retval, lasthash)

		log.Printf("sending block: %+v\n", newmsg)
		nc.Publish(msg.Reply, newmsg)
	case BLOCK_FOUND:
		Found = true
	}

	//	log.Printf("Msg: %+v\n", msg)

}
