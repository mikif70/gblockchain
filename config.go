// config.go
package main

const (
	MINE_RATE          = 1000
	INITIAL_DIFFICULTY = 18
)

var (
	GENESIS_DATA = Block{
		LastHash: []byte{00, 00},
		Data:     &Data{},
	}
)
