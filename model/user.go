package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var initializedInstance *Model = nil

func InitUser() *Model {
	if initializedInstance != nil {
		return initializedInstance
	}

	dir, _ := os.Getwd()
	userJsonPathHandler, _ := os.Open(dir + "/users.json")
	userData, _ := ioutil.ReadAll(userJsonPathHandler)
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(userData, &jsonData)

	if err != nil {
		return nil
	}

	initializedInstance = &Model{
		records: jsonData,
	}

	return initializedInstance
}
