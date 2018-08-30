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
		MakeExecutable() error
		MakeRPCCommand() (*RPCCommand, error)
		MakeAdminCommand() (*AdminCommand, error)
		ParseCommand() (string, string, error)
		ParseRPCParams() ([]*prot.OneCriteriaStruct, error)
		ParseAdminParams() ([]*prot.OneCriteriaStruct, error)
	}
)

func (in inputData) MakeExecutable() error {
	command, commandType, err := in.ParseCommand()
	if err != nil {
		return err
	}
	if commandType == "admin" && !in.user.IsAdmin() {
		return errors.New("Permission denied")
	}

	//bc := RPCCommand{command, commandType, params}
	return nil
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

func (in inputData) ParseRPCParams() ([]*prot.OneCriteriaStruct, error) {
	return nil, nil
}

func (in inputData) ParseAdminParams() ([]string, error) {
	return nil, nil
}
