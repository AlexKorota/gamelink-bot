package admincmd

import (
	"context"
	"gamelinkBot/command_list"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
	"strings"
)

type (
	//RqvokeFabric - struct for Revoke fabric
	RevokeFabric struct{}
	//RevokeCommand - struct for revoke command
	RevokeCommand struct {
		userName string
		params   []string
		res      iface.Responder
	}
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(RevokeFabric{})
}

//CommandName - return human readable command name
func (c RevokeFabric) CommandName() string {
	return command_list.CommandRevoke
}

//RequireAdmin - func for checking if admin permissions required
func (c RevokeFabric) RequireAdmin() bool {
	return true
}

//Require - return array of needed permissions
func (c RevokeFabric) Require() []string {
	return []string{command_list.CommandRevoke}
}

//TryParse - func for parsing request
func (c RevokeFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command RevokeCommand
		err     error
	)
	if command.userName, command.params, err = service.CompareParsePermissionCommand(req.Request(), "/"+command_list.CommandRevoke); err != nil {
		if err == service.UnknownCommandError {
			return nil, nil
		}
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (cc RevokeCommand) Execute(ctx context.Context) {
	user, err := Executor().RevokePermissions(cc.userName, cc.params)
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	if user == nil {
		cc.res.Respond("there is no admin with this name anymore")
	} else {
		cc.res.Respond("Success " + user.Name + " now has next permissions: " + strings.Join(user.Permissions, ", "))
	}
}
