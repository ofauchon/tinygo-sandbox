package core

import (
	"machine"

	"github.com/ofauchon/go-lorawan-stack"
	"tinygo.org/x/drivers/sx127x"
)

var (

	// Prepare for Lora Operation
	LoraConf = sx127x.LoraConfig{
		Freq:       868100000,
		Bw:         sx127x.SX127X_LORA_BW_125_0,
		Sf:         sx127x.SX127X_LORA_SF9,
		Cr:         sx127x.SX127X_LORA_CR_4_7,
		HeaderType: sx127x.SX127X_LORA_HEADER_EXPLICIT,
		Preamble:   12,
		//Ldr:            sx127x.SX127X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF,
		Iq:             sx127x.SX127X_LORA_IQ_STANDARD,
		Crc:            sx127x.SX127X_LORA_CRC_ON,
		SyncWord:       sx127x.SX127X_LORA_MAC_PUBLIC_SYNCWORD,
		LoraTxPowerDBm: 20,
	}

	UartConsole *machine.UART
	UartRS485   *machine.UART
	LoraStack   lorawan.LoraWanStack
	LoraRadio   *sx127x.Device
)

const (

	// Serial console
	UART1_TX_PIN = machine.PA9
	UART1_RX_PIN = machine.PA10

	// Serial to RS485
	UART2_TX_PIN = machine.PA2
	UART2_RX_PIN = machine.PA3

	// RFM95 SPI Connection to Bluepill
	SPI_SCK_PIN = machine.PA5
	SPI_SDO_PIN = machine.PA7
	SPI_SDI_PIN = machine.PA6
	SPI_CS_PIN  = machine.PB8
	SPI_RST_PIN = machine.PB9

	//RS485 module
	RS485_DERE_PIN = machine.PB0 // FIXME

	// DIO RFM95 Pin connection to BluePill
	DIO0_PIN        = machine.PA0
	DIO0_PIN_MODE   = machine.PinInputPullup
	DIO0_PIN_CHANGE = machine.PinRising

	DIO1_PIN        = machine.PA1
	DIO1_PIN_MODE   = machine.PinInputPullup
	DIO1_PIN_CHANGE = machine.PinRising
)
