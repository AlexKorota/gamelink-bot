package command

import (
	"gamelinkBot/prot"
	"github.com/kataras/iris/core/errors"
	"strings"
)

type (
	inputData struct {
		input string
		user  user
	}

	Maker interface {
		MakeExecutable() (*baseCommand, error)
		ParseCommand() (string, string, error)
		ParseParams() ([]*prot.OneCriteriaStruct, error)
	}
)

func (in inputData) MakeExecutable() (*baseCommand, error) {
	command, commandType, err := in.ParseCommand()
	if err != nil {
		return nil, err
	}
	if commandType == "admin" && !in.user.IsAdmin() {
		return nil, errors.New("Permission denied")
	}
	params, err := in.ParseParams()
	if err != nil {
		return nil, err
	}

	if commandType == "user" && !in.user.IsAdmin() {
		//проверка прав в базе данных
	}
	bc := baseCommand{command, commandType, params}
	return &bc, nil
}

func (in inputData) ParseCommand() (string, string, error) {
	commands := map[string]string{
		"/start":              "user",
		"/send_push":          "user",
		"/count":              "user",
		"/find":               "user",
		"/delete":             "user",
		"/update":             "user",
		"/get_user":           "user",
		"/grant_permissions":  "admin",
		"/revoke_permissions": "admin",
		"/show_permissions":   "admin",
	}
	command := in.input[:strings.Index(in.input, " ")]
	if _, ok := commands[command]; !ok {
		return "", "", errors.New("Invalid command")
	}
	return command, commands[command], nil
}

func (in inputData) ParseParams() ([]*prot.OneCriteriaStruct, error) {
	return nil, nil
}
