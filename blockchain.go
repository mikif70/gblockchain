package main

import (
	"encoding/json"
	"log"
)

type Chain []Block

func NewChain() *Chain {
	chain := Chain{}

	chain = append(chain, Genesis())

	return &chain
}

func (c *Chain) lastBlock() Block {
	return (*c)[len(*c)-1]
}

func (c *Chain) firstBlock() Block {
	return (*c)[0]
}

func (c *Chain) addBlock(data string) {
	if !ServerMode {

	}
	newBlock, _ := MineBlock(c.lastBlock(), []byte(data))
	*c = append(*c, newBlock)

}

func (c *Chain) replaceChain(chain Chain) {

	if ServerMode {
		if len(chain) <= len(*c) {
			return
		}

		if !isValidChain(&chain) {
			return
		}
	}

	(*c) = chain
}

func isValidChain(c *Chain) bool {
	genesis, _ := json.Marshal(Genesis())
	firstBlock, _ := json.Marshal(c.firstBlock())

	if string(genesis) != string(firstBlock) {
		log.Printf("genesis not found: first: %+v - genesis: %+v\n", string(firstBlock), string(genesis))
		return false
	}

	for i := 1; i < len(*c); i++ {
		block := (*c)[i]
		aLastHash := (*c)[i-1].Hash
		lastDifficulty := (*c)[i-1].Difficulty

		log.Printf("check %d: %s\n", i, string(block.Data))
		log.Printf("check %d - 1: %s - %x\n", i, string((*c)[i-1].Data), (*c)[i-1].LastHash[:4])

		if string(block.LastHash) != string(aLastHash) {
			log.Printf("invalid lasthash %d: %x != %x", i, block.LastHash[:4], aLastHash[:4])
			return false
		}

		validateHash := cryptoHash(&block)
		if string(block.Hash) != string(validateHash) {
			log.Printf("invalid hash: %x != %x", block.Hash[:4], validateHash[:4])
			return false
		}

		if (lastDifficulty - block.Difficulty) > 1 {
			log.Printf("invalid difficulty: %d != %d", lastDifficulty, block.Difficulty)
			return false
		}
	}

	return true
}
