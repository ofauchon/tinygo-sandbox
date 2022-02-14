package core

import (
	"errors"
	"machine"
	"strings"
	"time"
)

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

	case "version":
		println("version 1.0")

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
