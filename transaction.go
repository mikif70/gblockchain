// transaction.go
package main

import (
	"crypto/ecdsa"
	//	"time"
)

type Transaction struct {
	ID        string
	OutputMap *OutputMap
	Input     *Input
}

type OutputMap struct {
	Recipient    int
	SenderWallet *Wallet
}

type Input struct {
	Timestamp int64
	Amount    int
	Address   *ecdsa.PublicKey
	Signature string
}

/*
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
*/

/*
func createInput(senderWallet *Wallet, outputMap *OutputMap) *Input {
	input := Input{
		Timestamp: time.Now().Unix(),
		Amount:    senderWallet.Balance,
		Address:   senderWallet.PublicKey,
		Signature: senderWallet.sign(outputMap),
	}

	return &input
}
*/

/*
func createOutputMap(senderWallet *Wallet, recipient string, amount int) *OutputMap {

	omap := OutputMap{}

	omap.Recipient = amount
	omap.SenderWallet.PublicKey = senderWallet.Balance - amount

	return &omap
}

func verifySignature(publicKey string, data int, signature string) bool {

	return true
}
*/
