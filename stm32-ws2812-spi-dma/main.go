// Connects to an WS2812 RGB LED strip with DMA and SPI
//
// This application assumes your SPI bus is clocked at 8Mhz
// References:
// https://vivonomicon.com/2019/07/05/bare-metal-stm32-programming-part-9-dma-megamix/
// https://github.com/varunmadhavam/Tetris/blob/master/Core/Src/hw.c
// https://community.st.com/s/article/how-to-use-dmamux-on-the-stm32-microcontrollers

// Each led needs One R/G/B byte value
// Every ws2812 bit is encoded in a SPI Byte
// If we want to drive 16 LED, we need to send
// 16 * 3 * 8 SPI = 384 bytes (plus some additional RESET ZEROs)

package main

import (
	"crypto/rand"
	"device/stm32"
	"fmt"
	"machine"
	"time"
	"unsafe"
)

//var leds [10]color.RGBA

const (
	WS2812_NUM_LEDS    = 60
	WS2812_RESET_PULSE = 60
	WS2812_BUFFER_SIZE = WS2812_NUM_LEDS*24 + WS2812_RESET_PULSE
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

// GetRand16 returns 2 random bytes
func GetRand24() [3]uint8 {
	randomBytes := make([]byte, 3)
	_, err := rand.Read(randomBytes)
	if err == nil {
		r := [3]uint8{randomBytes[0], randomBytes[1], randomBytes[2]}
		return r
	} else {
		println("RAND ERROR ... ", stm32.RNG.SR.Get(), err)
		for {
		}
	}

}

// Get 3x8bits random array
func Rand24() [3]uint8 {
	res := randState
	res ^= res << 13
	res ^= res >> 17
	res ^= res << 15
	randState = res
	ret := [3]uint8{uint8(res & 0xFF), uint8((res >> 8) & 0xFF), uint8((res >> 16) & 0xFF)}
	return ret
}

// wsWriteColor will convert the value into a SPI bitstream
// We use 8 spi bits to encode a WS2812 pulse,
// So val byte will be written in a 8 byte array
func encodeColor(buffer []uint8, val uint8) {
	//println("encodeColor pointer:%p", &buffer[0], "lenbuf:", len(buffer))
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

// Changes RGB color of led at ledPos
func setLedColor(ledPos int, r, g, b uint8) {
	//println("Setcolor led", ledPos, "rgb:", r, g, b)
	start := ledPos * 24
	encodeColor(ws2812_buffer[start:], r)
	encodeColor(ws2812_buffer[start+8:], g)
	encodeColor(ws2812_buffer[start+16:], b)
}

// Enable  SPI TX DMA Circular transfert
func prepareDMASPI() {
	println("DMA.CCR1", stm32.DMA1.CCR1.Get())
	println("DMAMUX.C0CR", stm32.DMAMUX.C0CR.Get())

	// Enable DMA/DMAMUX clocks
	stm32.RCC.AHB1ENR.SetBits(stm32.RCC_AHB1ENR_DMAMUX1EN | stm32.RCC_AHB1ENR_DMA1EN | stm32.RCC_APB2ENR_SPI1EN)

	// Disable device, memory to memory mode, peripheral increment
	stm32.DMA1.CCR1.ClearBits(stm32.DMA_CCR1_MEM2MEM | stm32.DMA_CCR1_PINC | stm32.DMA_CCR1_EN)
	// Set Memory increment, circular mode, memory to peripheral
	//stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_MINC | stm32.DMA_CCR1_DIR)
	stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_MINC | stm32.DMA_CCR1_CIRC | stm32.DMA_CCR1_DIR)
	println("DMA.CCR1", stm32.DMA1.CCR1.Get())
	println("DMAMUX.C0CR", stm32.DMAMUX.C0CR.Get())

	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_MSIZE_Bits8, stm32.DMA_CCR1_MSIZE_Msk, 0) // 8bit mem mode
	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_PSIZE_Bits8, stm32.DMA_CCR1_PSIZE_Msk, 0) // 8bit periph mode
	stm32.DMA1.CCR1.ReplaceBits(stm32.DMA_CCR1_PL_High, stm32.DMA_CCR1_PL_Msk, 0)        // High priority

	stm32.DMAMUX.C0CR.ReplaceBits(stm32.DMAMUX_C0CR_DMAREQ_ID_SPI1_TX_DMA, stm32.DMAMUX_C0CR_DMAREQ_ID_Msk, 0) // Route DMA to SPI1 TX
	stm32.DMA1.CMAR1.Set(uint32(uintptr(unsafe.Pointer(&ws2812_buffer[0]))))                                   // Memory address
	stm32.DMA1.CPAR1.Set(uint32(uintptr(unsafe.Pointer(&machine.SPI0.Bus.DR))))                                // Periph. data register addr
	stm32.DMA1.CNDTR1.Set(WS2812_BUFFER_SIZE)
	println("Buf count:", stm32.DMA1.CNDTR1.Get())         // Transfert size
	fmt.Printf("Buf Addr: %X\r\n", stm32.DMA1.CMAR1.Get()) // Transfert size

	stm32.SPI1.CR1.ClearBits(stm32.SPI_CR1_SPE) // Disable SPI
	/*
		stm32.SPI1.CR1.ClearBits(stm32.SPI_CR1_LSBFIRST)                // MSB First
		stm32.SPI1.CR1.SetBits(stm32.SPI_CR1_SSM | stm32.SPI_CR1_SSI)   // Software management NSS
		stm32.SPI1.CR1.SetBits(stm32.SPI_CR1_MSTR)                      // Master
		stm32.SPI1.CR1.SetBits(stm32.SPI_CR1_CPOL | stm32.SPI_CR1_CPHA) // Clock polarity and phase
		stm32.SPI1.CR2.ReplaceBits(stm32.SPI_CR2_DS_EightBit, stm32.SPI_CR2_DS_Msk, 0)
		stm32.SPI1.CR1.ReplaceBits(stm32.SPI_CR1_BR_Div4, stm32.SPI_CR1_BR_Msk, 0) // Prescaler/8 '48Mh/8=8Mhz
	*/
	stm32.SPI1.CR2.SetBits(stm32.SPI_CR2_TXDMAEN) // Tx Buffer DMA
	stm32.DMA1.CCR1.SetBits(stm32.DMA_CCR1_EN)    // Enable DMA1 Channel 1
	stm32.SPI1.CR1.SetBits(stm32.SPI_CR1_SPE)     // Enable SPI

	println("DMA.CCR1", stm32.DMA1.CCR1.Get())
	println("DMAMUX.C0CR", stm32.DMAMUX.C0CR.Get())
	println("SPI1.CR1", stm32.SPI1.CR1.Get())
}

func clearAll() {
	for i := 0; i < WS2812_BUFFER_SIZE; i++ {
		ws2812_buffer[i] = 0x00
	}
}

func animate1() {

	dir := 1
	pos := 0
	oldpos := 0
	for {
		oldpos = pos
		if pos >= 60 && dir == 1 {
			dir = -1
		}
		if pos <= 1 && dir == -1 {
			dir = 1
		}

		pos = pos + dir
		println(pos)

		//clearAll()
		setLedColor(pos, 0xFF, 0xFF, 0xFF)
		setLedColor(oldpos, 0x00, 0x00, 0x00)

		time.Sleep(time.Millisecond * 250)

	}

}

func main() {
	println("WS2812 SPI DMA Transfert")

	println(machine.CPUFrequency())
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	},
	)

	wsInit()
	prepareDMASPI()

	/*
	 * I suspect a DMA buffer cache feature
	 * Sometimes changes in the buffer are not visible in SPI communications.
	 * Maybe we should switch OFF/ON DMA when changing the buffer
	 */

	/*
		loop := uint8(0)
		for {
			//r := Rand24()

			for i := 0; i < WS2812_NUM_LEDS; i++ {
				if (i % 3) == 0 {
					setLedColor(i, 0xFF, 0x00, 0x00)
				} else {
					setLedColor(i, 0x00, 0xFF, 0x00)
				}
			}

			loop++
			time.Sleep(time.Millisecond * 100)

		}
	*/

	animate1()

}
