package main

import (
	"strings"
	"errors"
)

const TX_TMOUT uint32 = 2000

func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("AT+FREQ: Set frequency")
		println("AT+RX: Receive mode until key pressed")
		println("AT+TX <msg>: Transmit msg string")

	case "AT+TX":
		if len(ss) == 2 {
			println("Lora TX message: ", ss[1])
			send_data := ss[1]
			err := loraRadio.LoraTx([]byte(send_data), TX_TMOUT)
			if err != nil {
				println("Send error", err)
			}
		}

	default:
		return errors.New("Unknown command")
	}

	return nil
}
