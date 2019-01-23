// protocol.go === Redis protocol
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/secmask/go-redisproto"
)

func handleConn(conn net.Conn, ui *gocui.Gui) {
	defer conn.Close()
	parser := redisproto.NewParser(conn)
	writer := redisproto.NewWriter(bufio.NewWriter(conn))
	var ew error
	for {
		command, err := parser.ReadCommand()
		if err != nil {
			_, ok := err.(*redisproto.ProtocolError)
			if ok {
				ew = writer.WriteError(err.Error())
			} else {
				log.Println(err, " closed connection to ", conn.RemoteAddr())
				SetNet(false, ui)
				break
			}
		} else {
			cmd := strings.ToUpper(string(command.Get(0)))
			switch cmd {
			case "BLOCKS":
				jchain, err := json.Marshal(chain)
				if err != nil {
					ew = writer.WriteBulkString(fmt.Sprintf("BLOCKS ERR: %+v\n", err))
				} else {
					ew = writer.WriteBulk(jchain)
				}
			case "MINE":
				data := &Data{}
				err := json.Unmarshal(command.Get(1), data)
				log.Printf("Data unmarshal: %+v\n", data)
				if err != nil {
					log.Printf("MINE ERR: %+v\n", err)
					ew = writer.WriteBulkString(fmt.Sprintf("MINE ERR: %+v\n", err))
					continue
				}
				chain.addBlock(*data)

				// PubSub synch
				// ....

				jchain, err := json.Marshal(chain)
				if err != nil {
					ew = writer.WriteBulkString(fmt.Sprintf("MINE ERR: %+v\n", err))
				} else {
					ew = writer.WriteBulk(jchain)
				}
			default:
				ew = writer.WriteError("Command not support")
			}
		}
		if command.IsLast() {
			writer.Flush()
		}
		if ew != nil {
			log.Println("Connection closed", ew)
			SetNet(false, ui)
			break
		}
	}
}

func NewListener(ui *gocui.Gui) {
	listener, err := net.Listen("tcp", ":5981")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error on accept: ", err)
			continue
		}
		SetNet(true, ui)
		go handleConn(conn, ui)
	}
}
