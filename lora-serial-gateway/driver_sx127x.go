//go:build sx127x

package main

import (
	"machine"
	"tinygo.org/x/drivers/sx127x"
)

var (
	loraRadio *sx127x.Device

	SX127X_PIN_RST  = machine.PB9
	SX127X_PIN_CS   = machine.PB8
	SX127X_PIN_DIO0 = machine.PA0
	SX127X_PIN_DIO1 = machine.PA1
	SX127X_SPI      = machine.SPI0
)

// Handle sx127x interrupts
func sx127xIntHandler() {
	loraRadio.HandleInterrupt()

}

func prepareRadio() {
	// Configure SX127x control GPIOS
	SX127X_PIN_RST.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SX127X_PIN_CS.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SX127X_PIN_DIO0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	SX127X_PIN_DIO1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	SX127X_SPI.Configure(machine.SPIConfig{Frequency: 500000, Mode: 0})

	println("Configure DIO interrupts")
	SX127X_PIN_DIO0.SetInterrupt(machine.PinRising|machine.PinFalling, func(machine.Pin) {
		sx127xIntHandler()
	})

	SX127X_PIN_DIO1.SetInterrupt(machine.PinRising|machine.PinFalling, func(machine.Pin) {
		sx127xIntHandler()
	})

	// SX127x driver init
	println("Init sx127x driver")
	loraRadio = sx127x.New(SX127X_SPI, SX127X_PIN_CS, SX127X_PIN_RST)

	// SX127x needs reset
	loraRadio.Reset()

	// Check the radio is ready
	println("Try to detect sx127x device")
	state := loraRadio.DetectDevice()

	if !state {
		panic("main: sx127x not detected... Aborting")
	}
	// Attach the Lora Radio to LoraStack
	loraStack.AttachLoraRadio(loraRadio)
}
