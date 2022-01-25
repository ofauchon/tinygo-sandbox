// Connects to an WS2812 RGB LED strip with 10 LEDS.
//
// See either the others.go or digispark.go files in this directory
// for the neopixels pin assignments.
package main

import (
	"device/stm32"
	"image/color"
	"machine"
	"time"
	"unsafe"
)

var leds [10]color.RGBA

const (
	WS2812_NUM_LEDS    = 3
	WS2812_RESET_PULSE = 60
	WS2812_BUFFER_SIZE = WS2812_NUM_LEDS*24 + WS2812_RESET_PULSE
)

var (
	ws2812_buffer []uint8
)

// WS2812 Initialization
func wsInit() {
	ws2812_buffer = make([]uint8, WS2812_BUFFER_SIZE, WS2812_BUFFER_SIZE)
	wsSendSpi()
}

// wsWriteColor will convert the value into a SPI bitstream
// We use 8 spi bits to encode a WS2812 pulse,
// So val byte will be written in a 8 byte array
func wsWriteColor(buffer []uint8, val uint8) {
	//	println("WriteColor val:", val, "lenbuf:", len(buffer))
	mask := uint8(0x80)
	for i := 7; i >= 0; i = i - 1 {
		mask = mask >> 1
		if (val & mask) > 0 {
			buffer[i] = 0xFC
		} else {
			buffer[i] = 0x80
		}
	}

}

func wsSendSpi() {
	err := machine.SPI0.Tx(ws2812_buffer, nil)
	if err != nil {
		panic(err)
	}
}

func wsSetColor(ledPos, r, g, b uint8) {
	println("Setcolor led", ledPos, "rgb:", r, g, b)
	wsWriteColor(ws2812_buffer[ledPos*24:], r)
	wsWriteColor(ws2812_buffer[ledPos*24+8:], g)
	wsWriteColor(ws2812_buffer[ledPos*24+16:], b)
}

// DMA configuration (channel 1).
// CCR register:
// - Memory-to-peripheral
// - Circular mode enabled.
// - Increment memory ptr, don't increment periph ptr.
// - -bit data size for both source and destination.
// - High priority.
func prepareDMA() {
	stm32.DMA1.CCR1.ClearBits(stm32.DMA_CCR1_MEM2MEM_Msk |
		stm32.DMA_CCR1_PL_Msk |
		stm32.DMA_CCR1_MSIZE_Msk |
		stm32.DMA_CCR1_PSIZE_Msk |
		stm32.DMA_CCR1_PINC_Msk |
		stm32.DMA_CCR1_EN_Msk)

	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_PL_High, stm32.DMA_CCR1_PL_Msk, 0) // High Priority

	stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_MINC_Msk | // Memory increment
		stm32.DMA_CCR1_CIRC_Msk | // Circular mode
		stm32.DMA_CCR1_DIR_Msk) // Memory to Peripheral

	// Route DMA to SPI1 TX
	stm32.DMAMUX.C0CR.ReplaceBits(stm32.DMAMUX_C0CR_DMAREQ_ID_SPI1_TX_DMA, stm32.DMAMUX_C0CR_DMAREQ_ID_Msk, 0)
	// Source: Address of the framebuffer.
	stm32.DMA1.CMAR1.Set(uint32(uintptr(unsafe.Pointer(&ws2812_buffer[0]))))
	// Destination: SPI1 data register.
	stm32.DMA1.CPAR1.Set(uint32(uintptr(unsafe.Pointer(&machine.SPI0.Bus.DR))))
	// Set DMA data transfer length (framebuffer length).
	stm32.DMA1.CNDTR1.Set(WS2812_BUFFER_SIZE)

}

func main() {
	//led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	println("HELLO")
	println(machine.CPUFrequency())
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		Mode:      0,
	},
	)

	wsInit()
	for {
		for i := uint8(0); i < WS2812_NUM_LEDS; i++ {
			wsSetColor(i, 0xAA, 0xAA, 0xAA)

			/*
				for k := 0; k < len(ws2812_buffer); k++ {
					print(ws2812_buffer[k], " ")
				}
			*/
			println("")
			machine.SPI0.Tx(ws2812_buffer, nil)
		}
		println("pause")
		time.Sleep(time.Second * 1)
	}

}
