// wallet.go
package main

import (
	"crypto/ecdsa"
)

type KeyPair struct {
	public  int
	private int
}

type Wallet struct {
	Balance   int
	KeyPair   *ecdsa.PrivateKey
	PublicKey *ecdsa.PublicKey
}

func NewWallet() *Wallet {

	keyp, _ := genKeyPair()

	wallet := &Wallet{
		Balance:   0,
		KeyPair:   keyp,
		PublicKey: &keyp.PublicKey,
	}

	return wallet
}

/*
func (w *Wallet) sign(data *Block) string {

	return w.KeyPair.sign(cryptoHash(data))
}
*/
