package util

import "math"

func Chunk(data []byte, size int) ([][]byte, int) {
	loops := int(math.Ceil(float64(len(data)) / float64(size)))
	chunks := make([][]byte, loops)
	for i := 0; i < loops; i++ {
		chunks[i] = data[(size * i):(size * (i + 1))]
	}
	return chunks, loops
}
