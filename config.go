// config.go
package main

import (
	"flag"
	//	"time"
)

const (
	MINE_RATE          = 10000
	INITIAL_DIFFICULTY = 10
	DEBUG              = false
	REFRESH_TIME       = 10
)

var (
	ServerMode    bool
	MinerMode     bool
	LogFilename   string
	MineRate      int
	servers       string
	ClientChannel = "blockchain"
	ServerChannel = "commands"
	MinerChannel  = "miner"
	ServerQueue   = "server"
)

func init() {
	flag.StringVar(&LogFilename, "l", "./chain.log", "log file")
	flag.BoolVar(&ServerMode, "s", false, "Server mode")
	flag.BoolVar(&MinerMode, "m", false, "miner mode")
	flag.IntVar(&MineRate, "t", MINE_RATE, "mine rate (msec)")
	flag.StringVar(&servers, "S", "nats://94.32.64.100:4222, nats://94.32.64.114:4222", "nats servers")
}
