package main

import (
	//	"math/rand"
	//	"encoding/json"
	"fmt"
	"log"

	"strings"

	"github.com/jroimartin/gocui"
	//ui "github.com/gizak/termui"
)

var (
	ui *gocui.Gui
)

func runConsole() {
	var err error

	ui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	//	ui.Highlight = true
	ui.BgColor = gocui.ColorBlack
	ui.FgColor = gocui.ColorWhite

	ui.SetManagerFunc(layout)

	ui.SetViewOnTop("data")

	keyBind(ui)
	mouseBind(ui)

	if err := ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keyBind(ui *gocui.Gui) {
	if err := ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func mouseBind(ui *gocui.Gui) {
	ui.Cursor = true
	ui.Mouse = true
	if err := ui.SetKeybinding("add", gocui.MouseLeft, gocui.ModNone, addBlock); err != nil {
		log.Panicln(err)
	}
	if err := ui.SetKeybinding("sync", gocui.MouseLeft, gocui.ModNone, synchBlockChain); err != nil {
		log.Panicln(err)
	}
	if err := ui.SetKeybinding("quit", gocui.MouseLeft, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := ui.SetKeybinding("data", gocui.MouseLeft, gocui.ModNone, setFocus); err != nil {
		log.Panicln(err)
	}
	if err := ui.SetKeybinding("data", gocui.MouseRelease, gocui.ModNone, setCursor); err != nil {
		log.Panicln(err)
	}
	if err := ui.SetKeybinding("data", gocui.MouseRight, gocui.ModNone, clear); err != nil {
		log.Panicln(err)
	}
}

func layout(ui *gocui.Gui) error {
	var difficulty, blocks, blockchain, data, net *gocui.View
	var addButton, idButton, quitButton *gocui.View
	var err error
	maxX, _ := ui.Size()

	if difficulty, err = ui.SetView("difficulty", 0, 0, 16, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		difficulty.Title = "Difficulty"
	}

	if blocks, err = ui.SetView("blocks", 0, 3, 16, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		blocks.Title = "Blocks"
	}

	if blocks, err = ui.SetView("nonce", 0, 6, 16, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		blocks.Title = "Last Nonce"
	}

	if blocks, err = ui.SetView("hash", 0, 9, 16, 11); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		blocks.Title = "Last Hash"
	}

	if blocks, err = ui.SetView("time", 0, 12, 16, 14); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		blocks.Title = "Time"
	}

	if net, err = ui.SetView("net", 0, 15, 12, 17); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		net.Title = "Net"
		if ServerMode {
			fmt.Fprint(net, " Server")
		} else if MinerMode {
			fmt.Fprint(net, " Miner")
		} else {
			fmt.Fprint(net, " Client")
		}
	}

	if blockchain, err = ui.SetView("blockchain", 18, 0, maxX-1, 18); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		blockchain.Autoscroll = true

		blockchain.Title = "timestamp  lasthash  hash    diff nounce   data"
		printBlockchain(ui)
	}

	if !MinerMode {
		if data, err = ui.SetView("data", 0, 19, maxX-16, 21); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			data.Title = "NewData"
			data.Editable = true
		}

		if addButton, err = ui.SetView("add", maxX-15, 19, maxX-1, 21); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			fmt.Fprintln(addButton, "     Add")
		}
	}

	if MinerMode {
		if idButton, err = ui.SetView("id", 0, 19, 42, 21); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			fmt.Fprintf(idButton, "id: %s\n", myID)
		}
	}

	if quitButton, err = ui.SetView("quit", maxX-15, 22, maxX-1, 24); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(quitButton, "    Quit")
	}

	return nil
}

func printBlockchain(ui *gocui.Gui) {
	v, err := ui.View("blockchain")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	if chain == nil {
		return
	}

	for i := range *chain {
		fmt.Fprintf(v, " %d %x %x %04d %06d %+v \n",
			(*chain)[i].Timestamp,
			string((*chain)[i].LastHash[0:4]),
			string((*chain)[i].Hash[0:4]),
			(*chain)[i].Difficulty,
			(*chain)[i].Nonce,
			string((*chain)[i].Data),
		)
	}
	printDifficulty(ui)
	printBlocksNum(ui)
	printNonce(ui)
	printLastHash(ui)
	printTime(ui)
}

func setNet(status bool) {
	v, err := ui.View("net")
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Status: %v - %s\n", status, v.Buffer())

	if status {
		v.Clear()
		v.BgColor = gocui.ColorYellow
		v.FgColor = gocui.ColorBlack
		fmt.Fprint(v, " Server")
	} else {
		v.Clear()
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite
		fmt.Fprint(v, " Server")
	}
}

func printTime(ui *gocui.Gui) {
	v, err := ui.View("time")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	fmt.Fprintf(v, " %d sec", ((*chain)[len(*chain)-1].Timestamp - (*chain)[len(*chain)-2].Timestamp))
}

func printLastHash(ui *gocui.Gui) {
	v, err := ui.View("hash")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	fmt.Fprintf(v, " %x", string((*chain)[len(*chain)-1].Hash[0:4]))
}

func printNonce(ui *gocui.Gui) {
	v, err := ui.View("nonce")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	fmt.Fprintf(v, " %d", (*chain)[len(*chain)-1].Nonce)
}

func printDifficulty(ui *gocui.Gui) {
	v, err := ui.View("difficulty")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	fmt.Fprintf(v, " %d", (*chain)[len(*chain)-1].Difficulty)
}

func printBlocksNum(ui *gocui.Gui) {
	v, err := ui.View("blocks")

	if err != nil {
		log.Panicln(err)
	}
	v.Clear()

	fmt.Fprintf(v, " %d", len(*chain))
}

func synchBlockChain(ui *gocui.Gui, v *gocui.View) error {
	printBlockchain(ui)

	return nil
}

func addBlock(ui *gocui.Gui, v *gocui.View) error {
	data, _ := ui.View("data")

	msg, err := makeMsg(ADD, []byte(strings.TrimSpace(data.Buffer())), []byte{})
	if err != nil {
		log.Printf("add error: %+v\n", err)
		return err
	}
	nc.Publish(ServerChannel, msg)
	data.Clear()
	data.SetCursor(0, 0)

	return nil
}

func deleteTemp(ui *gocui.Gui, v *gocui.View) error {
	return ui.DeleteView(v.Name())
}

func clear(ui *gocui.Gui, v *gocui.View) error {
	v.Clear()
	return nil
}

func quit(ui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCursor(ui *gocui.Gui, v *gocui.View) error {
	return v.SetCursor(0, 0)
}

func setFocus(ui *gocui.Gui, v *gocui.View) error {
	_, err := ui.SetCurrentView(v.Name())

	return err
}
