package util

import "encoding/base64"

func Byte2base64URI(data []byte) []byte {
	return []byte(
		"data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(
			data,
		),
	)
}
