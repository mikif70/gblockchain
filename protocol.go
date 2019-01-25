// protocol.go === Redis protocol
package main

import (
	//	"encoding/gob"
	"encoding/json"
	//	"log"
)

var (
	UPDATE      = "update"
	ADD         = "add"
	MINE        = "mine"
	INIT        = "init"
	BLOCK_FOUND = "found"
)

type myMsg struct {
	Cmd  string `json:cmd`
	Msg  []byte `json:msg`
	Data []byte `json:data`
}

func makeMsg(cmd string, data []byte, msg []byte) ([]byte, error) {
	mymsg := myMsg{
		Cmd:  cmd,
		Data: data,
		Msg:  msg,
	}

	jmsg, err := json.Marshal(mymsg)
	//	log.Printf("jmsg: %+v\n", string(jmsg))
	return jmsg, err
}
