//go:build dev
// +build dev

package main

import (
	"encoding/hex"
	"github.com/ofauchon/tinygo-sandbox/stm32-solivia-rs485-lorawan/core"
)

func LoraInitOTAA() {
	core.LoraStack.SetOtaa(
		[8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[8]uint8{0xA8, 0x40, 0x41, 0x00, 0x01, 0x81, 0xB3, 0x65},
		[16]uint8{0x2C, 0x44, 0xFC, 0xF8, 0x6C, 0x7B, 0x76, 0x7B, 0x8F, 0xD3, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16},
	)

	println("main: APPEUI:", hex.EncodeToString(core.LoraStack.Otaa.AppEUI[:]))
	println("main: DEVEUI:", hex.EncodeToString(core.LoraStack.Otaa.DevEUI[:]))
	println("main: APPKEY", hex.EncodeToString(core.LoraStack.Otaa.AppKey[:]))

}
