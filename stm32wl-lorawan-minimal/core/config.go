package core

import (
	"machine"

	"github.com/ofauchon/go-lorawan-stack"
	"tinygo.org/x/drivers/sx126x"
)

var (

	// Prepare for Lora Operation
	LoraConf = sx126x.LoraConfig{
		Freq:       868100000,
		Bw:         sx126x.SX126X_LORA_BW_125_0,
		Sf:         sx126x.SX126X_LORA_SF9,
		Cr:         sx126x.SX126X_LORA_CR_4_7,
		HeaderType: sx126x.SX126X_LORA_HEADER_EXPLICIT,
		Preamble:   12,
		//Ldr:            sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF,
		Iq:             sx126x.SX126X_LORA_IQ_STANDARD,
		Crc:            sx126x.SX126X_LORA_CRC_ON,
		SyncWord:       sx126x.SX126X_LORA_MAC_PUBLIC_SYNCWORD,
		LoraTxPowerDBm: 20,
	}

	UartConsole *machine.UART
	LoraStack   lorawan.LoraWanStack
	LoraRadio   *sx126x.Device
)

const (

	// Serial console
	UART1_TX_PIN = machine.PA9
	UART1_RX_PIN = machine.PA10
)
