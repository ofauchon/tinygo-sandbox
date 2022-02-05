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

// This will keep us connected
func loraConnect() {
	for {
		for !loraConnected {
			err := core.LoraStack.LoraWanJoin()
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

		println("cmdall")
		dat := sdec.GenCommand(0x01, [2]uint8{0x60, 0x01})
		core.RS485Send(dat)
		r := RS485Read(5) // Read for 5 seconds
		info, err := sdec.SoliviaParseInfoMsg(r)
		if err != nil {
			println(err)
			break
		}
		println("dbg/con: RAW:", hex.EncodeToString(info.LastPacket))
		println("dbg/con: ID:", info.Id, "PART:", info.PartNo, "SN:", info.SerialNo, "DATE:", info.DateCode)
		println("dbg/con: ACVolt:", info.ACVolt, "ACFreq:", info.ACFreq, "ACAmp:", info.ACAmp, "ACPower:", info.ACPower)
		println("dbg/con: DCVolt:", info.DCVolt, "DCAmp:", info.DCAmp)

		// Encode payload of Int/Ext sensors
		encoder.Reset()
		encoder.AddTemperature(1, float64(20)/1000)
		encoder.AddRelativeHumidity(2, float64(50)/100)
		encoder.AddTemperature(1, float64(10)/1000)
		encoder.AddRelativeHumidity(2, float64(80)/100)
		cayBytes := encoder.Bytes()

		if loraConnected {
			println("main: Sending LPP payload: ", hex.EncodeToString(cayBytes))
			err := core.LoraStack.LoraSendUplink(cayBytes)
			if err != nil {
				println(err)
			}
		}

		time.Sleep(time.Second * 180)
		machine.LED.Set(!machine.LED.Get())
	}
}
