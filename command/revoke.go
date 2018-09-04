package command

import (
	"context"
	"gamelinkBot/service"
)

type (
	RevokeFabric  struct{}
	RevokeCommand struct {
		userName string
		params   []string
		res      Responder
	}
)

const (
	commandRevoke = "grant_permissions"
)

func init() {
	SharedParser().RegisterFabric(RevokeFabric{})
}

func (c RevokeFabric) RequireAdmin() bool {
	return true
}

func (c RevokeFabric) Require() []string {
	return []string{commandRevoke}
}

func (c RevokeFabric) TryParse(req RequesterResponder) (Command, error) {
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

func (cc RevokeCommand) Execute(ctx context.Context) {
	//
	//if err != nil {
	//	cc.res.Respond(err.Error())
	//	return
	//}
	//
	//cc.res.Respond(r.String())
}
