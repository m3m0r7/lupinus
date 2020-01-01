package util

import (
	"math"
	"strings"
)

func Chunk(data []byte, size int) ([][]byte, int) {
	loops := int(math.Ceil(float64(len(data)) / float64(size)))
	chunks := make([][]byte, loops)
	for i := 0; i < loops; i++ {
		chunks[i] = data[(size * i):(size * (i + 1))]
	}
	return chunks, loops
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
