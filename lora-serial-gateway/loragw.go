package main

// Lora serial gateway

import (
	"encoding/hex"
	"machine"
	"time"

	"github.com/ofauchon/go-lorawan-stack"
	"tinygo.org/x/drivers/lora"
)



// serial() function is a goroutine for handling USART commands
func serial() string {
	input := make([]byte, 64) // serial port buffer
	i := 0

	for {

		if uart.Buffered() > 0 {

			data, _ := uart.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				uart.Write([]byte("\r\n"))
				cmd := string(input[:i])
				err := processCmd(cmd)
				if err != nil {
					println(err)
				}
				i = 0
			default: // pressed any other key
				uart.WriteByte(data)
				input[i] = data
				i++
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

}

func main() {
	println("***** TinyGo Lorawan Serial Gateway ****")

	// Configure LED GPIO
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := 0; i < 3; i++ {
		machine.LED.Low()
		time.Sleep(time.Millisecond * 250)
		machine.LED.High()
		time.Sleep(time.Millisecond * 250)
	}

	// This is specific to sx126x/sx127x (use -tags sx126x or sx127x )
	prepareRadio()

	// Initial lora configuration
	loraConf := lora.Config{
		Freq:           868100000,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       16,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPublic,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)




	// UART
	uart = machine.UART0
	uart.Configure(machine.UARTConfig{9600, 1, 0})
	go serial()


	for {
		machine.LED.Set(!machine.LED.Get())
		time.Sleep(1 * time.Second)
		cycle++
	}


}
