package util

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

func Byte2base64URI(data []byte) []byte {
	return []byte(
		"data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(
			data,
		),
	)
}

func Sha512(value string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(value)))
}

func Sha512WithSalt(value string, salt string) string {
	return Sha512(salt + value)
}
