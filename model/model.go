package model

import (
	"fmt"
	"../util"
)

type Model struct {
	Records map[string]interface{}
}

func (model *Model) Find(username string, password string) *map[string]interface{} {
	fmt.Printf("UserName: %s\n", username)
	fmt.Printf("Password: %s\n", password)

	ownProperty := util.GetFromMap(username, model.Records)
	if ownProperty == nil {
		return nil
	}

	path := model.Records[username].(map[string]interface{})

	passwordOwnProperty := util.GetFromMap("password", path)
	if passwordOwnProperty == nil  {
		return nil
	}

	// If password is correct return the user object
	if path["password"] == password {
		return &path
	}

	return nil
}

