package main

// Lorawan stack example code for sx127x

import (
	"encoding/hex"
	"machine"
	"time"

	"github.com/ofauchon/go-lorawan-stack"
	"tinygo.org/x/drivers/lora"
)

// Globals
var (
	loraStack     lorawan.LoraWanStack
	loraConnected bool
)

// This will keep us connected
func loraConnect() {
	for {
		for !loraConnected {
			err := loraStack.LoraWanJoin()
			if err != nil {
				println("main:Error joining Lorawan:", err, ", will wait 300 sec")
				time.Sleep(time.Second * 300)
			} else {
				println("main: Lorawan connection established")
				loraConnected = true
			}
		}
		// We are connected
		if loraConnected {
			machine.LED.Set(!machine.LED.Get())
		}
		time.Sleep(time.Second * 3)
	}
}

func main() {
	println("***** TinyGo Lorawan Stack Demo TOTO ****")

	// Configure LED GPIO
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := 0; i < 3; i++ {
		machine.LED.Low()
		time.Sleep(time.Millisecond * 250)
		machine.LED.High()
		time.Sleep(time.Millisecond * 250)
	}

	// Define OOTA settings
	// Temporary keys
	switch provider := "chirpstack"; provider {
	case "chirpstack":
		loraStack.SetOtaa(
			[8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[8]uint8{0xA8, 0x40, 0x41, 0x00, 0x01, 0x81, 0xB3, 0x65},
			[16]uint8{0x2C, 0x44, 0xFC, 0xF8, 0x6C, 0x7B, 0x76, 0x7B, 0x8F, 0xD3, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16},
		)
	case "ttn":
		loraStack.SetOtaa(
			[8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[8]uint8{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x04, 0xA9, 0x12},
			[16]uint8{0x67, 0x57, 0xBB, 0x98, 0x1D, 0x0E, 0x26, 0x71, 0xF4, 0x0F, 0x53, 0x4F, 0x6E, 0x4C, 0xD8, 0x7F},
		)
	case "orange":
		loraStack.SetOtaa(
			[8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[8]uint8{0x71, 0x33, 0x17, 0x88, 0x0C, 0x10, 0x88, 0x01},
			[16]uint8{0x61, 0x52, 0xB4, 0x33, 0x17, 0x12, 0x33, 0x44, 0xBE, 0xAF, 0xF0, 0x0F, 0x01, 0x02, 0x03, 0x01},
		)

	}
	println("main: APPEUI:", hex.EncodeToString(loraStack.Otaa.AppEUI[:]))
	println("main: DEVEUI:", hex.EncodeToString(loraStack.Otaa.DevEUI[:]))
	println("main: APPKEY", hex.EncodeToString(loraStack.Otaa.AppKey[:]))

	// This is specific to sx126x/sx127x (use -tags sx126x or sx127x )
	prepareRadio()

	// We need rand
	rnd := getRand16()
	loraStack.Otaa.DevNonce[0] = rnd[0]
	loraStack.Otaa.DevNonce[1] = rnd[1]
	println("DevNounce:", rnd[0], rnd[1])

	// Prepare for Lora Operation
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

	// Go routine for keeping us connected to Lorawan
	go loraConnect()

	// Wait 10 sec to give a chance to get a Lorawan connexion
	time.Sleep(time.Second * 20)

	payload := []byte("Hello from TinyGO")

	for {

		// Send payload if connected
		if loraConnected {
			println("main: Sending payload: ", hex.EncodeToString(payload))
			err := loraStack.LoraSendUplink(payload)
			if err != nil {
				println(err)
			}
		}

		// Sleep
		println("main: Sleep 180s")
		time.Sleep(180 * time.Second)

	}

}
