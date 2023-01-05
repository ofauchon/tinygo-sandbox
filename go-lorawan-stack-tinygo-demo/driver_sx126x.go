//go:build sx126x

package main

import (
	"device/stm32"
	rfswitch "github.com/ofauchon/go-lorawan-stack-tinygo-demo/rfswitch"
	"machine"
	"runtime/interrupt"
	"tinygo.org/x/drivers/sx126x"
)

var (
	loraRadio *sx126x.Device
)

// Handle sx126x interrupts
func radioIntHandler(intr interrupt.Interrupt) {
	loraRadio.HandleInterrupt()

}

func prepareRadio() {

	// Add interrupt handler for Radio IRQs (DIO)
	intr := interrupt.New(stm32.IRQ_Radio_IRQ_Busy, radioIntHandler)
	intr.Enable()

	// SX126x driver on SubGhzSPI (SPI3)
	loraRadio = sx126x.New(machine.SPI3)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Most boards have an RF FrontEnd Switch
	var radioSwitch rfswitch.CustomSwitch
	loraRadio.SetRfSwitch(radioSwitch)

	// Check the radio is ready
	state := loraRadio.DetectDevice()
	if !state {
		panic("main: sx126x not detected... Aborting")
	}

	// Attach the Lora Radio to LoraStack
	loraStack.AttachLoraRadio(loraRadio)
}
