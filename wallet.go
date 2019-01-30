// wallet.go
package main

type Wallet struct {
	Balance   int
	KeyPair   string
	PublicKey string
}

func NewWallet() *Wallet {

	keyp := genKeyPair()

	wallet := &Wallet{
		Balance:   0,
		KeyPair:   keyp,
		PublicKey: keyp.getPublic(),
	}

	return wallet
}

func (w *Wallet) sign(data interface{}) string {
	return w.KeyPair.sign(cryptoHash(data))
}
