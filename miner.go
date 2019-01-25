// miner.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/go-nats"
)

var (
	newblocks = make(map[string]bool)
)

func handleMinerMsg(msg *nats.Msg) {
	if msg == nil {
		log.Printf("msg err: %+v\n", msg)
		return
	}
	mymsg := myMsg{}
	err := json.Unmarshal(msg.Data, &mymsg)
	if err != nil {
		log.Printf("msg err: %+v\n", err)
		return
	}
	switch mymsg.Cmd {
	case MINE:
		go mineNewBlock(&mymsg, msg.Reply)
	case BLOCK_FOUND:
		lastHash := mymsg.Data
		newblocks[calcIdx(lastHash)] = true
		Found = true
	case UPDATE:
		newchain := Chain{}
		json.Unmarshal(mymsg.Data, &newchain)
		chain.replaceChain(newchain)
		printBlockchain(ui)
		log.Printf("Update: %x\n", newchain.lastBlock().Hash[:4])
	default:
	}
}

func mineNewBlock(mymsg *myMsg, reply string) {
	setNet(true)
	defer setNet(false)
	last := Block{}
	json.Unmarshal(mymsg.Msg, &last)
	id := calcIdx(last.LastHash[:8])
	newblocks[id] = false
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

	log.Printf("sending block: %x\n", minedBlock.Hash[:4])
	nc.Publish(reply, newmsg)
}

func MineBlock(last Block, data []byte) (Block, error) {
	var block Block
	var hash []byte
	var lastHash = last.Hash
	var difficulty = last.Difficulty
	var nonce = 0
	var timestamp int

	id := calcIdx(lastHash)

	for !isHashValid(hash, difficulty) {
		if !newblocks[id] {
			nonce += 1
			timestamp = int(time.Now().Unix())
			difficulty = AdjustDifficulty(&last, timestamp)
			block = NewBlock(int(timestamp), lastHash, []byte(data), nonce, difficulty)
			hash = cryptoHash(&block)
			log.Printf("Found: %t - %s: %d => %x\n", Found, data, nonce, hash[0:4])
			if DEBUG {
				log.Printf("hash: %d - %d - %+v\n", nonce, timestamp, hash)
			}
		} else {
			log.Printf("already found %x - mynonce %d ", hash[:4], nonce)
			return Block{}, errors.New(fmt.Sprintf("already found %x - mynonce %d ", hash[:4], nonce))
		}
	}

	block.Hash = hash

	return block, nil
}

func calcIdx(id []byte) string {
	return fmt.Sprintf("%x", id[:8])
}
