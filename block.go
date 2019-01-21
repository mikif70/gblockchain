package main

import (
	//	"bytes"
	//	"encoding/binary"
	//	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Data map[string]interface{}

type Block struct {
	Timestamp  int    `json:"timestamp"`
	LastHash   []byte `json:"lasthash"`
	Hash       []byte `json:"hash"`
	Nonce      int    `json:"nonce"`
	Difficulty int    `json:"difficulty"`
	Data       Data   `json:"data"`
}

func NewBlock(timestamp int, lastHash []byte, data Data, nonce int, difficulty int) Block {
	return Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       []byte{00, 00},
		Data:       data,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

func Genesis() Block {
	genesis := NewBlock(int(time.Now().Unix()), GENESIS_DATA.LastHash, GENESIS_DATA.Data, 0, INITIAL_DIFFICULTY)
	genesis.Hash = cryptoHash(&genesis)
	return genesis
}

func MineBlock(last Block, data Data) Block {
	var block Block
	var hash []byte
	var lastHash = last.Hash
	var difficulty = last.Difficulty
	var nonce = 0
	var timestamp int

	for !isHashValid(hash, difficulty) {
		nonce += 1
		timestamp = int(time.Now().Unix())
		difficulty = AdjustDifficulty(&last, timestamp)
		block = NewBlock(int(timestamp), lastHash, data, nonce, difficulty)
		hash = cryptoHash(&block)
		if DEBUG {
			fmt.Printf("hash: %d - %d - %+v\n", nonce, timestamp, hash)
		}
	}

	block.Hash = hash

	return block
}

func AdjustDifficulty(last *Block, timestamp int) int {

	var difference int

	difficulty := last.Difficulty

	if difficulty < 1 {
		if DEBUG {
			fmt.Printf("difference: 0\n")
		}
		return 1
	}

	difference = (timestamp - last.Timestamp) * 1000

	if DEBUG {
		fmt.Printf("difference %+v vs %d = %+v\n", difference, MINE_RATE, (difference - MINE_RATE))
	}

	if difference > MINE_RATE {
		return (difficulty - 1)
	}

	return (difficulty + 1)
}

func isHashValid(hash []byte, difficulty int) bool {

	if len(hash) == 0 {
		return false
	}

	zero := strings.Repeat("0", difficulty)
	hashDif := hash[0:difficulty]

	var hashZero, hashHex string

	for _, b := range hashDif {
		hashZero += fmt.Sprintf("%08b", b)
		if DEBUG {
			hashHex += fmt.Sprintf("%02x", b)
		}
	}

	if DEBUG {
		fmt.Printf("Zero: %+v - %+v - %+v\n", zero, hashZero, hashHex)
	}

	hashToDiff := hashZero[0:difficulty]

	if hashToDiff != zero {
		return false
	}

	return true
}
