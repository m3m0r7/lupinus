package subscriber

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"../validator"
	"../util"
	"os"
)

const (
	chunkSize = 8192
)

var (
	authKey = os.Getenv("AUTH_KEY")
	authKeySize = len(authKey)
)

func SubscribeImageStream(connection net.Conn) ([][]byte, int, error) {
	readAuthKey := make([]byte, authKeySize)
	receivedAuthKeySize, err := connection.Read(readAuthKey)
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		return nil, -1, err
	}

	// Compare the received auth key and settled auth key.
	if string(readAuthKey[:receivedAuthKeySize]) != authKey {
		fmt.Printf("err = %+v\n", err)
		return nil, -1, err
	}

	// Receive frame size
	frameSize := make([]byte, 4)
	_, errReceivingFrameSize := connection.Read(frameSize)
	if errReceivingFrameSize != nil {
		fmt.Printf("err = %+v\n", err)
		return nil, -1, err
	}

	realFrameSize := binary.BigEndian.Uint32(frameSize)
	realFrame := make([]byte, realFrameSize)

	receivedImageDataSize, errReceivingRealFrame := connection.Read(realFrame)
	if errReceivingRealFrame != nil {
		fmt.Printf("err = %+v\n", err)
		return nil, -1, err
	}

	frameData := realFrame[:receivedImageDataSize]

	if !validator.IsImageJpeg(frameData) {
		return nil, -1, errors.New("Does not match JPEG")
	}

	// Chunk the too long data.
	data, loops := util.Chunk(
		util.Byte2base64URI(
			frameData,
		),
		chunkSize,
	)
	return data, loops, nil
}

