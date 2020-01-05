package model

import (
	"encoding/json"
	"io/ioutil"
	"lupinus/config"
	"os"
)

var initializedInstance *Model = nil

func InitUser() *Model {
	if initializedInstance != nil {
		return initializedInstance
	}

	userJsonPathHandler, _ := os.Open(config.GetRootDir() + "/users.json")
	userData, _ := ioutil.ReadAll(userJsonPathHandler)
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(userData, &jsonData)

	if err != nil {
		return nil
	}

	initializedInstance = &Model{
		Records: jsonData,
	}

	return initializedInstance
}
