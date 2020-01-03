package subscriber

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lupinus/util"
	"lupinus/validator"
	"net"
	"os"
)

const (
	chunkSize = 8192
)

var (
	authKey = os.Getenv("AUTH_KEY")
	authKeySize = len(authKey)
)

func SubscribeImageStream(connection net.Conn) ([]byte, [][]byte, int, error) {
	readAuthKey := make([]byte, authKeySize)
	receivedAuthKeySize, err := connection.Read(readAuthKey)
	if err != nil {
		return nil, nil, -1, err
	}

	// Compare the received auth key and settled auth key.
	if string(readAuthKey[:receivedAuthKeySize]) != authKey {
		return nil, nil, -1, errors.New("Invalid auth key.")
	}

	// Receive frame size
	frameSize := make([]byte, 4)
	_, errReceivingFrameSize := connection.Read(frameSize)
	if errReceivingFrameSize != nil {
		return nil, nil, -1, errReceivingFrameSize
	}

	realFrameSize := binary.BigEndian.Uint32(frameSize)
	realFrame := []byte{}

	// Remaining calculator
	remaining := int(realFrameSize)
	for remaining > 0 {
		tmpRead := make([]byte, realFrameSize)
		receivedImageDataSize, errReceivingRealFrame := connection.Read(tmpRead)
		realFrame = append(realFrame, tmpRead...)
		if errReceivingRealFrame != nil {
			return nil, nil, -1, errReceivingRealFrame
		}

		remaining -= receivedImageDataSize
	}

	frameData := realFrame[:realFrameSize]

	if !validator.IsImageJpeg(frameData) {
		fmt.Printf("image = %d\n", realFrameSize)

		return nil, nil, -1, errors.New("Does not match JPEG")
	}

	// Chunk the too long data.
	data, loops := util.Chunk(
		util.Byte2base64URI(
			frameData,
		),
		chunkSize,
	)
	return frameData, data, loops, nil
}

