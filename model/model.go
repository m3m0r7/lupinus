package model

import (
	"../util"
	"os"
)

type Model struct {
	Records map[string]interface{}
}

func (model *Model) Find(username string, password string) *map[string]interface{} {

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
	if path["password"] == util.Sha512WithSalt(password, os.Getenv("SALT_KEY")) {
		return &path
	}

	return nil
}

