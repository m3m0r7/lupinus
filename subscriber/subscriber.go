package subscriber

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lupinus/util"
	"lupinus/validator"
	"net"
	"os"
	"strconv"
)

const (
	chunkSize = 8192
	protectedImageSize = 1024 * 1000 * 10
)

func SubscribeImageStream(connection net.Conn) ([]byte, [][]byte, int, error) {
	authKey := os.Getenv("AUTH_KEY")
	authKeySize := len(authKey)

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

	realFrameSize := int(binary.BigEndian.Uint32(frameSize))

	if realFrameSize < 0 || protectedImageSize < realFrameSize {
		return nil, nil, -1, errors.New(
			"protected memory allocation. tried to alloc = " + strconv.Itoa(realFrameSize),
		)
	}

	realFrame := []byte{}

	// Remaining calculator
	remaining := realFrameSize
	for remaining > 0 {
		tmpRead := make([]byte, remaining)
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

