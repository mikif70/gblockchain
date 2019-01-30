// transaction.go
package main

import (
	"time"
)

type Transaction struct {
	ID        string
	OutputMap *OutputMap
	Input     *Input
}

type OutputMap map[string]int

type Input struct {
	Timestamp int64
	Amount    int
	Address   string
	Signature string
}

func NewTransaction(senderWallet *Wallet, recipient string, amount int) *Transaction {
	id, err := newUUID()
	out := createOutputMap(senderWallet, recipient, amount)

	transaction := Transaction{
		ID:        id,
		OutputMap: out,
		Input:     createInput(senderWallet, out),
	}

	return &transaction
}

func createInput(senderWallet *Wallet, outputMap *OutputMap) *Input {
	input := Input{
		Timestamp: time.Now().Unix(),
		Amount:    senderWallet.Balance,
		Address:   senderWallet.PublicKey,
		Signature: senderWallet.sign(outputMap),
	}

	return &input
}

func createOutputMap(senderWallet *Wallet, recipient string, amount int) *OutputMap {

	omap := OutputMap{}

	omap[recipient] = amount
	omap[senderWallet.PublicKey] = senderWallet.Balance - amount

	return &omap
}

func verifySignature(publicKey string, data int, signature string) bool {

	return true
}
