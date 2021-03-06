package main

import (
	//	"bytes"
	//	"encoding/binary"
	//	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

var (
	Found        bool
	GENESIS_DATA = Block{
		LastHash:  []byte{00, 00, 00, 00, 00, 00, 00, 00},
		Data:      []byte("genesis"),
		Timestamp: int(time.Now().Unix()),
		Nonce:     0,
	}
)

type Block struct {
	Timestamp  int    `json:"timestamp"`
	LastHash   []byte `json:"lasthash"`
	Hash       []byte `json:"hash"`
	Nonce      int    `json:"nonce"`
	Difficulty int    `json:"difficulty"`
	Data       []byte `json:"data"`
}

func NewBlock(timestamp int, lastHash []byte, data []byte, nonce int, difficulty int) Block {
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
	genesis := NewBlock(GENESIS_DATA.Timestamp, GENESIS_DATA.LastHash, GENESIS_DATA.Data, GENESIS_DATA.Nonce, INITIAL_DIFFICULTY)
	genesis.Hash = cryptoHash(&genesis)
	return genesis
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
		fmt.Printf("difference %+v vs %d = %+v\n", difference, MineRate, (difference - MineRate))
	}

	if difference > MineRate {
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
