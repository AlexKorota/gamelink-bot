package admincmd

import (
	"context"
	"gamelinkBot/iface"
	"gamelinkBot/parser"
	"gamelinkBot/service"
	"strings"
)

type (
	//GrantFabric - struct for Grant fabric
	GrantFabric struct{}
	//GrantCommand - struct for grant command
	GrantCommand struct {
		userName string
		params   []string
		res      iface.Responder
	}
)

const (
	//commandGrant - const for command name
	commandGrant = "grant_permissions"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(GrantFabric{})
}

//CommandName - return human readable command name
func (c GrantFabric) CommandName() string {
	return commandGrant
}

//RequireAdmin - func for checking if admin permissions required
func (c GrantFabric) RequireAdmin() bool {
	return true
}

//Require - return array of needed permissions
func (c GrantFabric) Require() []string {
	return []string{commandGrant}
}

//TryParse - func for parsing request
func (c GrantFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command GrantCommand
		err     error
	)
	if command.userName, command.params, err = service.CompareParsePermissionCommand(req.Request(), "/"+commandGrant); err != nil {
		return nil, err
	}
	command.res = req
	return command, nil
}

//Execute - execute command
func (cc GrantCommand) Execute(ctx context.Context) {
	user, err := Executor().GrantPermissions(cc.userName, cc.params)
	if err != nil {
		cc.res.Respond(err.Error())
		return
	}
	cc.res.Respond("Success " + user.Name + " now has next permissions: " + strings.Join(user.Permissions, ", "))
}
