package core

import (
	"device/stm32"
	"machine"

	"runtime/interrupt"

	"tinygo.org/x/drivers/sx126x"

	rfswitch "github.com/ofauchon/tinygo-sandbox/stm32-join-lorawan/core/rfswitch"
)

// radioIntHandler will take care of radio interrupts
func radioIntHandler(intr interrupt.Interrupt) {
	LoraRadio.HandleInterrupt()
}

func InitRadio() {

	// Create the driver
	LoraRadio = sx126x.New(machine.SPI3)
	LoraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create RF Switch
	var radioSwitch rfswitch.CustomSwitch
	LoraRadio.SetRfSwitch(radioSwitch)

	// Detect the device
	state := LoraRadio.DetectDevice()
	if !state {
		panic("sx126x not detected.")
	}

	// Add interrupt handler for Radio IRQs
	intr := interrupt.New(stm32.IRQ_Radio_IRQ_Busy, radioIntHandler)
	intr.Enable()

	// Configure Lora settings (modulation, SF... etc )
	LoraRadio.LoraConfig(LoraConf)
}
