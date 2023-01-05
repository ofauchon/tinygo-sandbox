package main

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
			send_data = ss[1]
			err := d.Send([]byte(send_data))
			if err != nil {
				println("Send error", err)
			}
		}
	case "get":
		if len(ss) == 2 {
			switch ss[1] {
			case "freq":
				println("Freq:", d.GetFrequency())
			case "temp":
				temp, _ := d.ReadTemperature(0)
				println("Temperature:", temp)
			case "mode":
				mode := d.GetMode()
				println(" Mode:", mode)
			case "regs":
				for i := uint8(0); i < 0x60; i++ {
					val, _ := d.ReadReg(i)
					println(" Reg: ", strconv.FormatInt(int64(i), 16), " -> ", strconv.FormatInt(int64(val), 16))
				}
			default:
				return errors.New("Unknown command get")
			}
		}

	default:
		return errors.New("Unknown command")
	}

	return nil
}
