package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const TX_TMOUT uint32 = 2000
const RX_TMOUT uint32 = 2000

func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {

	case "?":
		println("AT+FREQ: Set frequency")
		println("AT+CONF: Display current configuration")
		println("AT+RX: Receive mode until key pressed")
		println("AT+TX <msg>: Transmit msg string")

	case "AT+CONF":
		fmt.Println("OK")
		fmt.Println(loraConf)

	case "AT+FREQ":
		if len(ss) == 2 {
			f, err := strconv.ParseUint(ss[1], 10, 32)
			if err == nil {
				loraConf.Freq = uint32(f)
				fmt.Println("OK")
			} else {
				return err
			}
		}

	case "AT+TX":
		if len(ss) == 2 {
			send_data := ss[1]
			err := loraRadio.LoraTx([]byte(send_data), TX_TMOUT)
			if err == nil {
				fmt.Println("OK")
			} else {
				return err
			}
		}

	case "AT+RX":
		d, err := loraRadio.LoraRx(RX_TMOUT)

	default:
		return errors.New("Unknown command")
	}

	return nil
}
