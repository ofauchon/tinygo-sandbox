// Receives data with LoRa.
package main

import (
	"encoding/hex"
	"machine"
	"time"

	cayennelpp "github.com/ofauchon/go-cayenne-lib"

	"github.com/ofauchon/tinygo-sandbox/stm32-solivia-rs485-lorawan/core"
)

// Globals
var (
	loraConnected bool
)

func ledMorse(val []uint16) {
	machine.LED.Low()
	for i := 0; i < len(val); i++ {
		machine.LED.High()
		time.Sleep(time.Millisecond * time.Duration(val[i]))
		machine.LED.Low()
		time.Sleep(time.Millisecond * 250)
	}
	machine.LED.Low()
}

// This will keep us connected
func loraConnect() {
	for {
		for !loraConnected {
			err := core.LoraStack.LoraWanJoin()
			if err != nil {
				println("loraConnect: Join error:", err)
				println("loraConnect: Wait 300 sec")
				time.Sleep(time.Second * 300)
			} else {
				println("loraConnect: Connected !")
				loraConnected = true
			}
		}
		// We are connected
		if loraConnected {
			ledMorse([]uint16{500, 2000})

		} else {
			ledMorse([]uint16{500, 500})
		}
		time.Sleep(time.Second * 10)
	}
}

func main() {

	// LED Init
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(true)

	// Console Init on UART1
	core.ConsoleInit(machine.UART1)
	core.ConsoleStartTask()

	println("Delta Solivia RS485 Lora gateway")

	// Init RS485
	println("main: Init RS485")
	core.RS485Init(machine.UART2)

	// Lora OTAA (Keys).
	// Should be declared in config_xxxx.go
	println("main: Load OTAA Keys")
	LoraInitOTAA()
	println("main: APPEUI:", hex.EncodeToString(core.LoraStack.Otaa.AppEUI[:]))
	println("main: DEVEUI:", hex.EncodeToString(core.LoraStack.Otaa.DevEUI[:]))
	println("main: APPKEY", hex.EncodeToString(core.LoraStack.Otaa.AppKey[:]))

	println("main: Init Radio Module")
	core.InitRadio()

	println("main: Get Rand Uint32 from LoraStack")
	rnd := core.LoraRadio.RandomU32()
	core.LoraStack.Otaa.DevNonce[0] = uint8((rnd) & 0xFF)
	core.LoraStack.Otaa.DevNonce[1] = uint8((rnd >> 8) & 0xFF)

	println("main: Attach Radio to Lora Stack")
	core.LoraStack.AttachLoraRadio(core.LoraRadio)

	println("main: Start Lora Stack loop")
	go loraConnect()

	// Wait 10 sec to give a chance to get a Lorawan connexion
	time.Sleep(time.Second * 20)

	// We'll encode with Cayenne LPP protocol
	encoder := cayennelpp.NewEncoder()

	var sdec core.SoliviaDecoder

	// Loop forever
	for {

		println("main: Get infos from PV")
		dat := sdec.GenCommand(0x01, [2]uint8{0x60, 0x01})
		core.RS485Send(dat)

		r := core.RS485Read(5) // Read for 5 seconds
		info, err := sdec.SoliviaParseInfoMsg(r)
		if err != nil {
			ledMorse([]uint16{500, 250, 250, 250, 250, 250, 250})
			println("main: Communication error:", err)
			println("main: Waiting for 60 seconds")
			time.Sleep(time.Second * 60)
			continue
		}
		println("main: Solivia infos:")
		println("      RAW:", hex.EncodeToString(info.LastPacket))
		println("      ID:", info.Id, "PART:", info.PartNo, "SN:", info.SerialNo, "DATE:", info.DateCode)
		println("      ACVolt:", info.ACVolt, "ACFreq:", info.ACFreq, "ACAmp:", info.ACAmp, "ACPower:", info.ACPower)
		println("      DCVolt:", info.DCVolt, "DCAmp:", info.DCAmp)

		// TESTS : info := &core.SoliviaInfos{ACAmp: 30, ACVolt: 2200, ACFreq: 499, DCVolt: 200, DCAmp: 30}

		// Encode payload of Int/Ext sensors
		encoder.Reset()
		encoder.AddVoltage(1, float64(info.ACVolt))
		encoder.AddFrequency(1, float64(info.ACFreq))
		encoder.AddCurrent(1, float64(info.ACAmp))
		encoder.AddPower(1, float64(info.ACPower))

		encoder.AddVoltage(2, float64(info.DCVolt))
		encoder.AddCurrent(2, float64(info.DCAmp))

		cayBytes := encoder.Bytes()

		if loraConnected {
			println("main: Sending LPP payload: ", hex.EncodeToString(cayBytes))
			ledMorse([]uint16{500, 250, 250, 250})
			err := core.LoraStack.LoraSendUplink(cayBytes)
			if err != nil {
				println(err)
			}
		}

		time.Sleep(time.Second * 180)
		machine.LED.Set(!machine.LED.Get())
	}
}
