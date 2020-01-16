package helper

import (
	"encoding/json"
	"fmt"
	"lupinus/config"
	"os"
	"path/filepath"
	"time"
)

func CreateStaticImage(image []byte, filename string) {
	path := config.GetRootDir() + "/storage/" + filename

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		fmt.Println(err)
		return
	}
	handle, _ := os.Create(path)

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
