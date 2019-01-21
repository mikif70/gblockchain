package main

import (
	"encoding/json"
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

func (c *Chain) addBlock(data Data) {
	newBlock := MineBlock(c.lastBlock(), data)
	*c = append(*c, newBlock)
}

func (c *Chain) replaceChain(chain Chain) {
	if len(chain) <= len(*c) {
		return
	}

	if !isValidChain(&chain) {
		return
	}

	(*c) = chain
}

func isValidChain(c *Chain) bool {
	genesis, _ := json.Marshal(Genesis())
	firstBlock, _ := json.Marshal(c.firstBlock())

	if string(genesis) != string(firstBlock) {
		return false
	}

	for i := 1; i < len(*c); i++ {
		block := (*c)[i]
		aLastHash := (*c)[i-1].LastHash
		lastDifficulty := (*c)[i-1].Difficulty

		if string(block.LastHash) != string(aLastHash) {
			return false
		}

		validateHash := cryptoHash(&block)
		if string(block.Hash) != string(validateHash) {
			return false
		}

		if (lastDifficulty - block.Difficulty) > 1 {
			return false
		}
	}

	return true
}
