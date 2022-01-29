// Receives data with LoRa.
package main

import (
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/ofauchon/tinygo-sandbox/stm32-solivia-rs485-lorawan/core"
)

// Globals
var (
	loraConnected bool
)

// Handle sx127x interrupts
func radioIntHandler(intr interrupt.Interrupt) {
	core.LoraRadio.HandleInterrupt()

}

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
	core.RS485Init(machine.UART2)

	// Lora OTAA (Keys).
	// Should be declared in config_xxxx.go
	LoraInitOTAA()

	// Initialize DevNonce
	rnd := core.LoraRadio.RandomU32()
	println("main: Rand32:", rnd)
	core.LoraStack.Otaa.DevNonce[0] = uint8((rnd) & 0xFF)
	core.LoraStack.Otaa.DevNonce[1] = uint8((rnd >> 8) & 0xFF)

	// Configure SPI, Lora SX127x
	core.InitLora()

	// Attach the Lora Radio to LoraStack
	core.LoraStack.AttachLoraRadio(core.LoraRadio)

	// Go routine for keeping us connected to Lorawan
	go loraConnect()

	// Wait 10 sec to give a chance to get a Lorawan connexion
	time.Sleep(time.Second * 20)

	// We'll encode with Cayenne LPP protocol
	//	encoder := cayennelpp.NewEncoder()

	// Loop forever
	for {
		time.Sleep(time.Second)
		machine.LED.Set(!machine.LED.Get())
	}
}
