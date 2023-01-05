//go:build !(bluepill && sx127x)

package main

// GetRand16 uses board-specific RNG from tinygo

import "crypto/rand"

func getRand16() [2]uint8 {
	randomBytes := make([]byte, 2)
	_, err := rand.Read(randomBytes)
	if err == nil {
		r16 := [2]uint8{randomBytes[0], randomBytes[1]}
		return r16
	} else {
		println("RAND ERROR ... Abort")
		for {
		}
	}

}
