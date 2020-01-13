package util

import (
	"strings"
)

func Chunk(data []byte, size int) ([][]byte, int) {
	chunks := make([][]byte, 1)
	chunks[0] = data
	return chunks, 1
	// TODO: FIX bottom code
	//	loops := int(math.Floor(float64(len(data)) / float64(size)))
	//	chunks := make([][]byte, loops)
	//	for i := 0; i < loops; i++ {
	//		chunks[i] = data[(size * i):(size * (i + 1))]
	//	}
	//
	//	lastPos := loops * size
	//	// Read remaining data
	//	remainBytes := len(data) - lastPos
	//	if remainBytes > 0 {
	//		// Expand slices
	//		newChunks := make([][]byte, loops + 1)
	//		newChunks = append(
	//			newChunks,
	//			chunks...,
	//		)
	//		newChunks[loops] = data[lastPos:(lastPos + remainBytes)]
	//	}
	//	return chunks, loops
}

func SplitWithFiltered(data string, sep string) []string {
	newSeparated := []string{}
	separated := strings.Split(data, sep)

	for _, val := range separated {
		if val == "" {
			continue
		}
		newSeparated = append(newSeparated, val)
	}
	return newSeparated
}

func GetFromMap(key string, mapObject map[string]interface{}) interface{} {
	val, ok := mapObject[key]
	if !ok {
		return nil
	}
	return val
}
