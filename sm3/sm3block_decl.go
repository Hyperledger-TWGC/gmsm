// +build amd64,386

package sm3

import (
	"encoding/binary"
)

// func block(msg []byte, p []byte)

func FF0(v1, v2, v3 int32) int32
func FF1(v1, v2, v3 int32) int32

func GG0(v1, v2, v3 int32) int32
func GG1(v1, v2, v3 int32) int32

func P0(x uint32) uint32
func P1(x uint32) uint32

// go:noescape
func block(dig *digest, b []byte)

// func mod(a int) int

// func TestInt(a uint32, b []byte)

type digest struct {
	h  [8]uint32
	ns int
}

func Write(b []byte) []byte {
	dig := &digest{
		h: [8]uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e},
	}
	b = pad(b)
	block(dig, b)
	result := make([]byte, 32)
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(result[i*4:], dig.h[i])
	}
	return result
}

func pad(msg []byte) []byte {
	length := len(msg) * 8
	msg = append(msg, 0x80) // Append '1'
	blockSize := 64         // Append until the resulting message length (in bits) is congruent to 448 (mod 512)
	for len(msg)%blockSize != 56 {
		msg = append(msg, 0x00)
	}
	// append message length
	msg = append(msg, uint8(length>>56&0xff))
	msg = append(msg, uint8(length>>48&0xff))
	msg = append(msg, uint8(length>>40&0xff))
	msg = append(msg, uint8(length>>32&0xff))
	msg = append(msg, uint8(length>>24&0xff))
	msg = append(msg, uint8(length>>16&0xff))
	msg = append(msg, uint8(length>>8&0xff))
	msg = append(msg, uint8(length>>0&0xff))

	if len(msg)%64 != 0 {
		panic("------SM3 Pad: error msgLen =")
	}
	return msg
}

func leftRotate(x uint32, i uint32) uint32 { return (x<<(i%32) | x>>(32-i%32)) }

func TestValue() uint32 {
	w16 := uint32(0x61626364)
	w9 := uint32(0x61626364)
	w3 := leftRotate(0xa121a024, 15)
	w13 := leftRotate(0x61626364, 7)
	w6 := uint32(0x61626364)
	return P1(w16^w9^w3) ^ w13 ^ w6
}

func p1Generic(x uint32) uint32 {
	return x ^ leftRotate(x, 15) ^ leftRotate(x, 23)
}
