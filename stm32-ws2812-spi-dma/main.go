// Connects to an WS2812 RGB LED strip with DMA and SPI
//
// This application assumes your SPI bus is configured and clocked at 8Mhz
// References:
// https://vivonomicon.com/2019/07/05/bare-metal-stm32-programming-part-9-dma-megamix/
// https://github.com/varunmadhavam/Tetris/blob/master/Core/Src/hw.c
// https://community.st.com/s/article/how-to-use-dmamux-on-the-stm32-microcontrollers
// https://www.st.com/resource/en/product_training/STM32G0-System-Direct-memory-access-controller-DMA.pdf
//
// Each led needs One R/G/B byte value
// Every ws2812 bit is encoded in a SPI Byte
// If we want to drive 16 LED, we need to send
// 16 * 3 * 8 SPI = 384 bytes (plus some additional RESET ZEROs)

package main

import (
	"device/stm32"
	"fmt"
	"machine"
	"time"
	"unsafe"
)

const (
	WS2812_NUM_LEDS    = 60
	WS2812_POST_ZEROS  = 60
	WS2812_PRE_ZEROS   = 2
	WS2812_BUFFER_SIZE = WS2812_PRE_ZEROS + WS2812_NUM_LEDS*24 + WS2812_POST_ZEROS
	DEBUG              = false
)

var (
	ws2812_buffer []uint8
	randState     = uint32(0x6542)
	loop          = int(0)
)

// WS2812 Initialization
func wsInit() {
	ws2812_buffer = make([]uint8, WS2812_BUFFER_SIZE, WS2812_BUFFER_SIZE)
}

// Prepare for SPI TX DMA Transfert
func prepareTransfert() {
	stm32.RCC.AHB1ENR.SetBits(stm32.RCC_AHB1ENR_DMAMUX1EN | stm32.RCC_AHB1ENR_DMA1EN | stm32.RCC_APB2ENR_SPI1EN) // Enable clocks

	stm32.DMA1.CCR1.ClearBits(stm32.DMA_CCR1_EN | stm32.DMA_CCR1_MEM2MEM | stm32.DMA_CCR1_PINC)                // Disable device and wrong config.
	stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_MINC | stm32.DMA_CCR1_DIR)                                          // Set Memory increment,  memory to peripheral
	stm32.DMAMUX.C0CR.ReplaceBits(stm32.DMAMUX_C0CR_DMAREQ_ID_SPI1_TX_DMA, stm32.DMAMUX_C0CR_DMAREQ_ID_Msk, 0) // Route DMA to SPI1 TX

	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_MSIZE_Bits8, stm32.DMA_CCR1_MSIZE_Msk, 0) // 8bit mem mode
	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_PSIZE_Bits8, stm32.DMA_CCR1_PSIZE_Msk, 0) // 8bit periph mode
	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_PL_High, stm32.DMA_CCR1_PL_Msk, 0)        // High priority

	stm32.SPI1.CR1.ClearBits(stm32.SPI_CR1_SPE)   // Disable SPI
	stm32.SPI1.CR2.SetBits(stm32.SPI_CR2_TXDMAEN) // Tx Buffer DMA
	stm32.SPI1.CR1.SetBits(stm32.SPI_CR1_SPE)     // Enable SPI
}

// Starts a SPI DMA Xfer
func startTransfert() {
	// Clear previous IRQ
	isrFlag := stm32.DMA1.ISR.Get()
	stm32.DMA1.IFCR.Set(isrFlag)

	stm32.DMA1.CCR1.ClearBits(stm32.DMA_CCR1_EN)     // Stop DMA to reconfigure it
	for stm32.DMA1.CCR1.HasBits(stm32.DMA_CCR1_EN) { // Wait it's actualy stopped
	}
	stm32.DMA1.CMAR1.Set(uint32(uintptr(unsafe.Pointer(&ws2812_buffer[0]))))    // Memory address
	stm32.DMA1.CPAR1.Set(uint32(uintptr(unsafe.Pointer(&machine.SPI0.Bus.DR)))) // Periph. data register addr
	stm32.DMA1.CNDTR1.Set(WS2812_BUFFER_SIZE)                                   // Transfer size
	stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_EN)                                  // Re-enable DMA

}

// "Kind of" Random Color picker
func Rand24() [3]uint8 {
	res := randState
	res ^= res << 13
	res ^= res >> 17
	res ^= res << 15
	randState = res
	ret := [3]uint8{uint8(res & 0xFF), uint8((res >> 8) & 0xFF), uint8((res >> 16) & 0xFF)}
	return ret
}

// encodeColor will convert a color value into a SPI bitstream
// We encode each  WS2812 pulses (0 or 1) with one SPI byte.
// So "val" byte will be written in a 8 byte array
func encodeColor(buffer []uint8, val uint8) {
	if DEBUG {
		println("encodeColor pointer:%p", &buffer[0], "lenbuf:", len(buffer))
	}
	mask := uint8(0x80)
	for i := 7; i >= 0; i = i - 1 {
		if (val & mask) > 0 {
			buffer[i] = 0xF0
		} else {
			buffer[i] = 0xC0
		}
		mask = mask >> 1
	}
}

// setLedColor Sets RGB color at a specific LED position
func setLedColor(ledPos int, r, g, b uint8) {
	start := WS2812_PRE_ZEROS + ledPos*24
	encodeColor(ws2812_buffer[start:], r)
	encodeColor(ws2812_buffer[start+8:], g)
	encodeColor(ws2812_buffer[start+16:], b)
}

// Clears all colors
func clearAll() {
	for i := 0; i < WS2812_BUFFER_SIZE; i++ {
		ws2812_buffer[i] = 0x00
	}
	for i := 0; i < WS2812_NUM_LEDS; i++ {
		setLedColor(i, 0x00, 0x00, 0x00)
	}
}

func dumpBuf() {
	println("")
	for k := 0; k < len(ws2812_buffer); k++ {
		if k%8 == 0 {
			fmt.Printf("\r\n> [%p] ", &ws2812_buffer[k])
		}
		print(ws2812_buffer[k], " ")
	}
	println("")
}

func animate1() {
	pos := int8(0) // Start at position zero
	curCol := Rand24()

	for {
		setLedColor(int(pos), curCol[0], curCol[1], curCol[2])
		startTransfert() // Start DMA XFer

		if pos > WS2812_NUM_LEDS {
			pos = 0
			loop++
			curCol = Rand24()
		} else {
			pos++
		}

		time.Sleep(time.Millisecond * 10)
	}

}

//
//
//
func main() {
	println("WS2812 SPI DMA Transfert")

	println(machine.CPUFrequency())
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	},
	)

	wsInit()
	prepareTransfert()

	// Start animation
	animate1()

}
