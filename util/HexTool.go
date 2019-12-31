package util

import "encoding/binary"

func Uint2bytes(i int, size int) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	return bytes[8-size : 8]
}