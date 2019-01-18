// config.go
package main

const (
	MINE_RATE          = 1000
	INITIAL_DIFFICULTY = 22
	DEBUG              = true
)

var (
	GENESIS_DATA = Block{
		LastHash: []byte{00, 00},
		Data:     &Data{},
	}
)
