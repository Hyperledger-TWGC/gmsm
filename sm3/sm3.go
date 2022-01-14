// +build amd64,386

package sm3

import "hash"

type SM3 struct{
	digest      [8]uint32 // digest represents the partial evaluation of V
	length      uint64    // length of the message
	unhandleMsg []byte    // uint8  //
}

// func Block(msg []byte) []byte {

// }

func goWrite(h [8]uint32, b []byte) {
	var w [68]uint32
	// w[0] = binary.BigEndian.Uint32(b[:4])
	h[7] = w[0]
}

func Sm3Sum(b []byte) []byte {
	return Write(b)
}

// 创建哈希计算实例
func New() hash.Hash {
	var sm3 SM3

	sm3.Reset()
	return &sm3
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (sm3 *SM3) BlockSize() int { return 64 }

// Size returns the number of bytes Sum will return.
func (sm3 *SM3) Size() int { return 32 }

// Reset clears the internal state by zeroing bytes in the state buffer.
// This can be skipped for a newly-created hash state; the default zero-allocated state is correct.
func (sm3 *SM3) Reset() {
	// Reset digest
	sm3.digest[0] = 0x7380166f
	sm3.digest[1] = 0x4914b2b9
	sm3.digest[2] = 0x172442d7
	sm3.digest[3] = 0xda8a0600
	sm3.digest[4] = 0xa96f30bc
	sm3.digest[5] = 0x163138aa
	sm3.digest[6] = 0xe38dee4d
	sm3.digest[7] = 0xb0fb0e4e

	sm3.length = 0 // Reset numberic states
	sm3.unhandleMsg = []byte{}
}

// Write (via the embedded io.Writer interface) adds more data to the running hash.
// It never returns an error.
func (sm3 *SM3) Write(p []byte) (int, error) {
	toWrite := len(p)
	sm3.length += uint64(len(p) * 8)
	msg := append(sm3.unhandleMsg, p...)
	// Update unhandleMsg
	sm3.unhandleMsg = msg[:]
	return toWrite, nil
}

// 返回SM3哈希算法摘要值
// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (sm3 *SM3) Sum(in []byte) []byte {
	_, _ = sm3.Write(in)
	return Write(sm3.unhandleMsg)
}