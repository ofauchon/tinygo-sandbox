//go:build bluepill && sx127x

package main

// GetRand16 gets random from SX126x radio module (RSSI value)

func getRand16() [2]uint8 {
	i := loraRadio.RandomU32()
	return [2]uint8{uint8(i & 0xFF), uint8(((i >> 8) & 0xFF))}

}
