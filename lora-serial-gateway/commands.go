package main

import (
	"errors"
	"strings"
)

const TX_TMOUT uint32 = 2000
const RX_TMOUT uint32 = 10000

func atoi(str string) (uint32, error) {
	var num uint32
	for _, c := range str {
		digit := int(c - '0')
		if digit < 0 || digit > 9 {
			return 0, errors.New("Invalid number")
		}
		num = num*10 + uint32(digit)
	}
	return num, nil
}

func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {

	case "?", "h", "help":

		println("AT+CONF: Display Lora configu")
		println("AT+FREQ: Set frequency")
		println("AT+BW: Set bandwidth")
		println("AT+SF: Set spreading factor")
		println("AT+CR: Set coding rate")
		println("AT+IQ: Set I/Q inversion")
		println("AT+CRC Set CRC enable")
		println("AT+HEADERTYPE: Set header type")
		println("AT+SYNCWORD: Set sync word")
		println("AT+TX: Send packet")
		println("AT+RX: Receive packet")

	case "AT+CONF":
		println("OK")
		println("Frequency:", loraConf.Freq, "Hz")
		println("Bandwidth:", loraConf.Bw)
		println("Spreading factor:", loraConf.Sf)
		println("Coding rate:", loraConf.Cr)
		println("Header type:", loraConf.HeaderType)
		println("Preamble length:", loraConf.Preamble)
		println("ldr:", loraConf.Ldr)
		println("IQ inversion:", loraConf.Iq)
		println("CRC:", loraConf.Crc)
		println("Sync word:", loraConf.SyncWord)
		println("Tx Power :", loraConf.LoraTxPowerDBm, "DBm")

	case "AT+FREQ":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Freq = uint32(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+BW":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Bw = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}
	case "AT+SF":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Sf = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+CR":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Cr = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+IQ":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Iq = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+CRC":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.Crc = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+HEADERTYPE":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.HeaderType = uint8(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+SYNCWORD":
		if len(ss) == 2 {
			f, err := atoi(ss[1])
			if err == nil {
				loraConf.SyncWord = uint16(f)
				loraRadio.LoraConfig(loraConf)
				println("OK")
			} else {
				return err
			}
		}

	case "AT+TX":
		if len(ss) == 2 {
			send_data := ss[1]
			err := loraRadio.LoraTx([]byte(send_data), TX_TMOUT)
			if err == nil {
				println("OK")
			} else {
				return err
			}
		}

	case "AT+RX":
		println("OK")
		r, err := loraRadio.LoraRx(RX_TMOUT)
		println("RX TIMEOUT")
		if err != nil {
			println("RX:", string(r))
		} else {
			return err
		}

	default:
		return errors.New("Unknown command")
	}

	return nil
}
