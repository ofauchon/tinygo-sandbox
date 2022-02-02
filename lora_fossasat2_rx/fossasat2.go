package main

// This application uses a SX1262 433Mhz to listen to Fossat2 Lora Frames
// (LoRa 401.7 Mhz SF: 11 CR: 8 BW: 125 kHz)

import (
	"device/stm32"
	"encoding/hex"
	"machine"
	"runtime/interrupt"

	rfswitch "tinygo.org/x/drivers/examples/sx126x/rfswitch"

	"tinygo.org/x/drivers/sx126x"
)

const FREQ = 868100000

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

var (
	loraRadio *sx126x.Device
	txmsg     = []byte("Hello TinyGO")
)

// radioIntHandler will take care of radio interrupts
func radioIntHandler(intr interrupt.Interrupt) {
	loraRadio.HandleInterrupt()
}

func main() {
	println("\n# TinyGo FossaSAT2 Receiver")
	println("# ---------------------------")
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Create the driver
	loraRadio = sx126x.New(machine.SPI3)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create RF Switch
	var radioSwitch rfswitch.CustomSwitch
	loraRadio.SetRfSwitch(radioSwitch)

	// Detect the device
	state := loraRadio.DetectDevice()
	if !state {
		panic("sx126x not detected.")
	}

	// Add interrupt handler for Radio IRQs
	intr := interrupt.New(stm32.IRQ_Radio_IRQ_Busy, radioIntHandler)
	intr.Enable()

	loraConf := sx126x.LoraConfig{
		Freq:           401700000,
		Bw:             sx126x.SX126X_LORA_BW_125_0,
		Sf:             sx126x.SX126X_LORA_SF11,
		Cr:             sx126x.SX126X_LORA_CR_4_8,
		HeaderType:     sx126x.SX126X_LORA_HEADER_EXPLICIT, //?
		Preamble:       12,                                 //?
		Ldr:            sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF,
		Iq:             sx126x.SX126X_LORA_IQ_STANDARD,
		Crc:            sx126x.SX126X_LORA_CRC_ON,
		SyncWord:       sx126x.SX126X_LORA_MAC_PRIVATE_SYNCWORD,
		LoraTxPowerDBm: 20,
	}
	// (LoRa 401.7 Mhz SF: 11 CR: 8 BW: 125 kHz)

	loraRadio.LoraConfig(loraConf)

	for {
		//tStart := time.Now()

		for {

			buf, err := loraRadio.LoraRx(LORA_DEFAULT_RXTIMEOUT_MS)

			if err != nil {
				println("RX Error: ", err)
			} else if buf != nil {
				println("Packet Received: len=", len(buf), " data:", hex.EncodeToString(buf))

			}
		}

	}

}
