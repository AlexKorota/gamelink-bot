package generalcmd

import (
	"context"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
)

type (
	//DeleteFabric - struct for DeleteFabric
	DeleteFabric struct{}
	//DeleteCommand - struct for delete command
	DeleteCommand struct {
		params []*iface.OneCriteriaStruct
		res    iface.Responder
	}
)

const (
	//commandDelete - const for command
	commandDelete = "delete"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(DeleteFabric{})
}

//CommandName - return human readable command name
func (c DeleteFabric) CommandName() string {
	return commandDelete
}

//RequireAdmin - func for checking if admin permissions required
func (c DeleteFabric) RequireAdmin() bool {
	return false
}

//Require - return array of needed permissions
func (c DeleteFabric) Require() []string {
	return []string{commandDelete}
}

//TryParse - func for parsing request
func (c DeleteFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command DeleteCommand
		err     error
	)
	if command.params, err = service.CompareParseCommand(req.Request(), "/"+commandDelete); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (cc DeleteCommand) Execute(ctx context.Context) {
	r, err := Executor().Delete(ctx, cc.params)
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond(r.String())
}
