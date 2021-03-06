package core

import (
	"encoding/hex"
	"errors"
	"machine"
	"strings"
	"time"
)

var sdec SoliviaDecoder

//ConsoleInit() Initializes console's UART
func ConsoleInit(port *machine.UART) {
	UartConsole = machine.UART1
	UartConsole.Configure(machine.UARTConfig{TX: UART1_TX_PIN, RX: UART1_RX_PIN, BaudRate: 115200})
}

//ConsoleStartTask() starts the console goroutine
func ConsoleStartTask() {
	go consoleTask()
}

//processCmd() handles console inputs commmand
func processCmd(cmd string) error {
	println("Console command:", cmd)

	ss := strings.Split(cmd, " ")
	if len(ss) == 0 {
		return errors.New("Bad command")
	}
	switch ss[0] {
	case "send": // send 02050102600185FC03
		if len(ss) == 2 {
			dat, err := DecodeString(ss[1])
			if err == nil {
				RS485Send(dat)
				rsp := RS485Read(3)
				println("dbg/con: Eead ", len(rsp))
			} else {
				println("Error in send command: ", err)
			}

		}

	case "cmd": // ex : CMD 1001
		if len(ss) == 2 {
			c, err := DecodeString(ss[1])
			if err == nil && len(c) == 2 {
				println("dbg/con: Send command:", ss[1])
				dat := sdec.GenCommand(0x01, [2]uint8{c[0], c[1]})
				RS485Send(dat)
				rsp := RS485Read(3)
				print("dbg/con: RS485 Receive: ")
				for _, v := range rsp {
					print(v, " ")
				}
				println("")

			}

		}

	case "dereon":
		println("DERE PIN ON")
		RS485_DERE_PIN.Set(true)

	case "dereoff":
		println("DERE PIN OFF")
		RS485_DERE_PIN.Set(false)

	case "cmdall":
		println("cmdall")
		dat := sdec.GenCommand(0x01, [2]uint8{0x60, 0x01})
		RS485Send(dat)
		r := RS485Read(5)
		info, err := sdec.SoliviaParseInfoMsg(r)
		if err != nil {
			println(err)
			break
		}
		println("dbg/con: RAW:", hex.EncodeToString(info.LastPacket))
		println("dbg/con: ID:", info.Id, "PART:", info.PartNo, "SN:", info.SerialNo, "DATE:", info.DateCode)
		println("dbg/con: ACVolt:", info.ACVolt, "ACFreq:", info.ACFreq, "ACAmp:", info.ACAmp, "ACPower:", info.ACPower)
		println("dbg/con: DCVolt:", info.DCVolt, "DCAmp:", info.DCAmp)

	default:
		println("Error processing command:", cmd)
		println("Usage:")
		println("cmdall : Request all datas")
		println("cmd 1001 : Send command 0x1001")
		println("cmd 1001 : Send command 0x1001")
		println("dereon : Switch DE/RE pin ON ")
		println("dereoff : Switch DE/RE pin OFF ")
	}

	return nil
}

// consoleTask() is the real go console routine
func consoleTask() string {

	println("Starting console task.")
	inputConsole := make([]byte, 128) // serial port buffer

	for {

		// Process console messages
		for UartConsole.Buffered() > 0 {

			data, _ := UartConsole.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				UartConsole.Write([]byte("\r\n"))
				cmd := string(inputConsole)
				err := processCmd(cmd)
				if err != nil {
					println(err)
				}
				inputConsole = nil
			default: // pressed any other key
				UartConsole.WriteByte(data)
				inputConsole = append(inputConsole, data)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

}
