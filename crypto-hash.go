package main

import (
	"crypto/sha256"
	//	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG = true
)

func cryptoHash(block *Block) []byte {
	start := time.Now()

	data, err := json.Marshal(block.Data)
	if err != nil {
		fmt.Printf("Error converting the Data: %+v\n", err)
		data = []byte{}
	}

	dataSlice := []string{
		strconv.Itoa(block.Timestamp),
		string(block.LastHash),
		strconv.Itoa(block.Nonce),
		strconv.Itoa(block.Difficulty),
		string(data),
	}

	sort.Strings(dataSlice)

	str := strings.Join(dataSlice, "")

	var hash = sha256.New()
	hash.Write([]byte(str))
	hash.Write([]byte(hash.Sum(nil)))

	if DEBUG {
		fmt.Printf("time: %+v\n", time.Since(start))
	}

	return hash.Sum(nil)
}
