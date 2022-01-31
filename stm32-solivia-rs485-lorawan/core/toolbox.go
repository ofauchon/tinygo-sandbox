package core

// This is the placeholder for various functions.

import (
	"errors"
)

var ErrLength = errors.New("encoding/hex: odd length hex string")

func EncodedLen(n int) int { return n * 2 }

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	// We can use the source slice itself as the destination
	// because the decode loop increments by one and then the 'seen' byte is not used anymore.
	n, err := Decode(src, src)
	return src[:n], err
}

func Decode(dst, src []byte) (int, error) {
	i, j := 0, 1
	for ; j < len(src); j += 2 {
		a, ok := fromHexChar(src[j-1])
		if !ok {
			return i, errors.New("Hex Decode err 1")
		}
		b, ok := fromHexChar(src[j])
		if !ok {
			return i, errors.New("Hex Decode err 2")
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if _, ok := fromHexChar(src[j-1]); !ok {
			return i, errors.New("Hex Decode err 3")
		}
		return i, ErrLength
	}
	return i, nil
}
