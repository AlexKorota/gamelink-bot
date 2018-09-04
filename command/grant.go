package command

import (
	"context"
	"gamelinkBot/service"
)

type (
	GrantFabric  struct{}
	GrantCommand struct {
		userName string
		params   []string
		res      Responder
	}
)

const (
	commandGrant = "grant_permissions"
)

func init() {
	SharedParser().RegisterFabric(GrantFabric{})
}

func (c GrantFabric) RequireAdmin() bool {
	return true
}

func (c GrantFabric) Require() []string {
	return []string{commandGrant}
}

func (c GrantFabric) TryParse(req RequesterResponder) (Command, error) {
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

func (cc GrantCommand) Execute(ctx context.Context) {
	//if err != nil {
	//	cc.res.Respond(err.Error())
	//	return
	//}
	//cc.res.Respond(r.String())
}
