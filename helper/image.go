package helper

import (
	"encoding/json"
	"lupinus/config"
	"os"
	"time"
)

func CreateStaticImage(image []byte, filename string) {
	handle, _ := os.Create(config.GetRootDir() + "/storage/" + filename)

	// Write simple image
	handle.Write(image)

	handle.Close()

	// Write meta
	metaHandle, _ := os.Create(config.GetRootDir() + "/storage/" + filename + ".meta.json")

	jsonData, _ := json.Marshal(
		map[string]interface{}{
			// FIXME: get extension by image data
			"extension":     "jpg",
			"time":          time.Now().Unix(),
			"camera_number": 0,
		},
	)

	metaHandle.Write(jsonData)

	metaHandle.Close()
}
