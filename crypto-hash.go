package main

import (
	"crypto/sha256"
	//	"encoding/hex"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
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

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
