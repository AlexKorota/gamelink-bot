package admincmd

import (
	"context"
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

const (
	//commandRevoke - const for revoke command name
	commandRevoke = "revoke_permissions"
)

//init - func for register fabric in parser
func init() {
	parser.SharedParser().RegisterFabric(RevokeFabric{})
}

//CommandName - return human readable command name
func (c RevokeFabric) CommandName() string {
	return commandRevoke
}

//RequireAdmin - func for checking if admin permissions required
func (c RevokeFabric) RequireAdmin() bool {
	return true
}

//Require - return array of needed permissions
func (c RevokeFabric) Require() []string {
	return []string{commandRevoke}
}

//TryParse - func for parsing request
func (c RevokeFabric) TryParse(req iface.RequesterResponder) (iface.Command, error) {
	var (
		command RevokeCommand
		err     error
	)
	if command.userName, command.params, err = service.CompareParsePermissionCommand(req.Request(), "/"+commandRevoke); err != nil {
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
	cc.res.Respond("Success " + user.Name + " now has next permissions: " + strings.Join(user.Permissions, ", "))
}
